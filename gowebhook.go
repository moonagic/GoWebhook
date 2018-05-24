package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	targetSecret     = ""
	targetURL        = ""
	targetPort       = ""
	targetHost       = ""
	targetShell      = ""
	targetLogDir     = "/etc/gowebhook/"
	targetLogFile    = "log"
	queue []string
	running          = false
)

func generateHashSignature(message string) string {
	h := hmac.New(sha1.New, []byte(targetSecret))
	h.Write([]byte(message))
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	log2file(string(r.URL.Host))
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func autoBuild(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "post" || r.Method == "POST") && r.URL.RequestURI() == targetURL {
		if r.Header.Get("x-github-event") == "push" {
			bodyContent, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			signature := r.Header.Get("X-Hub-Signature")
			if verifySignature(signature, string(bodyContent)) {
				fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
				log2file("验证通过,启动部署任务")
				queue = append(queue, "1")
				checkoutTaskStatus()
			} else {
				log2file("验证失败")
				fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
			}
		} else {
			fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch x-github-event\"}")
		}
	} else {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method or unknow request url\"}")
	}
}

func loadConfig() {
	result, err := ioutil.ReadFile("/etc/gowebhook/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localUrl, ok0 := m["requestUrl"].(string)
		localSecret, ok1 := m["secret"].(string)
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localShell, ok4 := m["script"].(string)
		if ok0 && ok1 && ok2 && ok3 && ok4 {
			targetURL = localUrl
			targetSecret = localSecret
			targetHost = localHost
			targetPort = localPort
			targetShell = localShell
		} else {
			log2file("Broken config.")
			os.Exit(0)
		}
	} else {
		log2file("Can not find config file...in \"/etc/gowebhook/config\"")
		os.Exit(0)
	}
}

func checkoutTaskStatus()  {
	if running {
		return
	}
	if len(queue) > 0 {
		data := []string{""}
		queue = data[:0:0]
		go startTask()
	}
}

func startService() {
	http.HandleFunc("/", index)
	http.HandleFunc("/auto_build", autoBuild)


	log2file(fmt.Sprintf("service starting... %s:%s", targetHost, targetPort))
	listenErr := http.ListenAndServe(fmt.Sprintf("%s:%s", targetHost, targetPort), nil)
	if listenErr != nil {
		log.Fatal("ListenAndServe: ", listenErr)
	}
}

func startTask() {
	running = true
	cmd := exec.Command("/bin/sh", targetShell)
	_, err := cmd.Output()
	if err == nil {
		running = false
		log2file("部署成功")
		checkoutTaskStatus()
	} else {
		running = false
		log2file(fmt.Sprintf("部署失败:\n %s", err))
		checkoutTaskStatus()
	}
}

func verifySignature(signature string, data string) bool {
	return signature == generateHashSignature(string(data))
}

func log2file(content string)  {
	var err error

	if _, err := os.Stat(targetLogDir); err == nil {
		fmt.Println("Dir exists", targetLogDir)
	} else {
		fmt.Println("Dir not exists, try to create...", targetLogDir)
		err := os.MkdirAll(targetLogDir, 0711)
		if err != nil {
			fmt.Println("Error creating directory", targetLogDir)
			fmt.Println("err:", err)
			return
		}
	}

	if _, err := os.Stat(targetLogDir + targetLogFile); err == nil {
		fmt.Println("Path exists", targetLogDir + targetLogFile)
	} else {
		fmt.Println("Path not exists, try to create...", targetLogDir + targetLogFile)
		_, err := os.Create(targetLogDir + targetLogFile)
		if err != nil {
			fmt.Println("Error creating file", targetLogDir + targetLogFile)
			fmt.Println("err:", err)
			return
		}
	}

	f, err := os.OpenFile(targetLogDir + targetLogFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err == nil {
		f.WriteString(content)
		f.WriteString("\n")
	} else {
		fmt.Println("Open file faild...", targetLogDir + targetLogFile)
		fmt.Println("err:", err)
	}
}

func main() {
	loadConfig()
	startService()
}
