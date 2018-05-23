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

const (
	version = 0.1
)

var (
	targetSecret = ""
	targetURL    = ""
	targetPort   = ""
	targetHost   = ""
	targetShell  = ""
)

func generateHashSignature(message string) string {
	h := hmac.New(sha1.New, []byte(targetSecret))
	h.Write([]byte(message))
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(string(r.URL.Host))
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
				fmt.Println("验证通过,启动部署任务")
				go startTask()
			} else {
				fmt.Println("验证失败")
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
			fmt.Println("Broken config.")
			os.Exit(0)
		}
	} else {
		fmt.Println("Can not find config file...in \"/etc/gowebhook/config\"")
		os.Exit(0)
	}
}

func startService() {
	http.HandleFunc("/", index)
	http.HandleFunc("/auto_build", autoBuild)

	fmt.Println("service starting...", targetHost, targetPort)
	listenErr := http.ListenAndServe(fmt.Sprintf("%s:%s", targetHost, targetPort), nil)
	if listenErr != nil {
		log.Fatal("ListenAndServe: ", listenErr)
	}
}

func startTask() {
	cmd := exec.Command("/bin/sh", targetShell)
	_, err := cmd.Output()
	if err == nil {
		fmt.Println("部署成功")
	} else {
		fmt.Println("部署失败:", err)
	}
}

func verifySignature(signature string, data string) bool {
	return signature == generateHashSignature(string(data))
}

func main() {
	loadConfig()
	startService()
}
