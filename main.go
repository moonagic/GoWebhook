package main

import (
	"net/http"
	"fmt"
	"crypto/sha1"
	"crypto/hmac"
	"encoding/hex"
	"io/ioutil"
	"os/exec"
	"encoding/json"
	"os"
	"log"
)

var TARGETSECRET string = ""
var TARGETURL string = ""
var TARGETPORT string = ""
var TARGETHOST string = ""
var TARGETSHELL string = ""


func generateHashSignature(message string) string {
	h := hmac.New(sha1.New, []byte(TARGETSECRET))
	h.Write([]byte(message))
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func autoBuild(w http.ResponseWriter, r *http.Request)  {
	if (r.Method == "post" || r.Method == "POST") && r.URL.RequestURI() == TARGETURL {
		fmt.Println(r.URL.RequestURI())
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

func loadConfig()  {
	result, err := ioutil.ReadFile("/etc/gowebhook/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localUrl, ok0 := m["requestUrl"].(string)
		LocalSecret, ok1 := m["secret"].(string)
		LocalHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localShell, ok4 := m["script"].(string)
		if ok0 && ok1 && ok2 && ok3 && ok4 {
			TARGETURL = localUrl
			TARGETSECRET = LocalSecret
			TARGETHOST = LocalHost
			TARGETPORT = localPort
			TARGETSHELL = localShell
		} else {
			fmt.Println("Broken config.")
			os.Exit(0)
		}
	} else {
		fmt.Println("Can not find config file...in \"/etc/gowebhook/config\"")
		os.Exit(0)
	}
}

func startService()  {
	http.HandleFunc("/", index)
	http.HandleFunc("/auto_build", autoBuild)

	fmt.Println("service starting...", TARGETHOST, TARGETPORT)
	listenErr := http.ListenAndServe(fmt.Sprintf("%s:%s", TARGETHOST, TARGETPORT), nil)
	if listenErr != nil {
		log.Fatal("ListenAndServe: ", listenErr)
	}
}

func startTask()  {
	cmd := exec.Command("/bin/sh", TARGETSHELL)
	_, err := cmd.Output()
	if err == nil {
		fmt.Println("部署成功")
	} else {
		fmt.Println("部署失败:", err)
	}
}

func verifySignature(signature string, data string) bool {
	return (signature == generateHashSignature(string(data)))
}

func main() {
	loadConfig()
	startService()
}
