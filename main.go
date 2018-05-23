package main

import (
	"log"
	"net/http"
	"fmt"
	"crypto/sha1"
	"crypto/hmac"
	"encoding/hex"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"os"
)

var TARGETSECRET string = ""
var TARGETURL string = ""
var TARGETPORT string = ""
var TARGETHOST string = ""
var TARGETSHELL string = ""


func ComputeHmac1(message string) string {
	h := hmac.New(sha1.New, []byte(TARGETSECRET))
	h.Write([]byte(message))
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func aotoBuild(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "post" || r.Method == "POST" {
		fmt.Println(r.URL.RequestURI())
		if r.URL.RequestURI() == TARGETURL {
			if r.Header.Get("x-github-event") == "push" {
				result, _ := ioutil.ReadAll(r.Body)
				r.Body.Close()

				signature := r.Header.Get("X-Hub-Signature")
				if signature == ComputeHmac1(string(result)) {
					fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
					fmt.Println("验证通过")
					go func() {
						cmd := exec.Command("/bin/sh", TARGETSHELL)
						bytes, err := cmd.Output()
						if err == nil {
							fmt.Println(string(bytes))
						} else {
							fmt.Println("err", err)
						}
					}()

				} else {
					fmt.Println("验证失败")
					fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
				}
			} else {
				fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch x-github-event\"}")
			}
		} else {
			fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unknow request url\"}")
		}
	} else {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method\"}")
	}
}

func main() {
	// config
	result, err := ioutil.ReadFile("/etc/webhook/config")
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
			fmt.Println("Broken config .")
			os.Exit(0)
		}
	} else {
		fmt.Println("Can not find config file.../etc/webhook/config")
		os.Exit(0)
	}
	
	// router
	http.HandleFunc("/", index)
	http.HandleFunc("/auto_build", aotoBuild)


	fmt.Println("service starting...", TARGETHOST, TARGETPORT)
	listenErr := http.ListenAndServe(fmt.Sprintf("%s:%s", TARGETHOST, TARGETPORT), nil)
	if listenErr != nil {
		log.Fatal("ListenAndServe: ", listenErr)
	}
}
