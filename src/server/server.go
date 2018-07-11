package server

import (
	"GoWebhooks/src/config"
	"GoWebhooks/src/task"
	"GoWebhooks/src/utils"
	"fmt"
	"io/ioutil"
	"net/http"
)

// StartService start service
func StartService(address string, port string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/auto_build", autoBuild)

	utils.Log2file(fmt.Sprintf("service starting... %s:%s", address, port))
	return http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	utils.Log2file(string(r.URL.Host))
	fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func autoBuild(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "post" || r.Method == "POST") && r.URL.RequestURI() == config.GetURL() {
		if r.Header.Get("x-github-event") == "push" {
			bodyContent, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			signature := r.Header.Get("X-Hub-Signature")
			if utils.VerifySignature(signature, string(bodyContent), config.GetSecret()) {
				fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
				utils.Log2file("验证通过,启动部署任务")
				task.AddNewTask(string(bodyContent))
			} else {
				utils.Log2file("验证失败")
				fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
			}
		} else {
			fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch x-github-event\"}")
		}
	} else {
		fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method or unknow request url\"}")
	}
}
