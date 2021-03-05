package server

import (
	"GoWebhook/src/config"
	"GoWebhook/src/task"
	"GoWebhook/src/utils"
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
	utils.Log2file(r.URL.Host)
	_, _ = fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func autoBuild(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "post" || r.Method == "POST") && r.URL.RequestURI() == config.GetURL() {
		if r.Header.Get("x-github-event") == "push" {
			bodyContent, _ := ioutil.ReadAll(r.Body)
			_ = r.Body.Close()
			signature := r.Header.Get("X-Hub-Signature")
			if utils.VerifySignature(signature, string(bodyContent), config.GetSecret()) {
				_, _ = fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
				utils.Log2file("验证通过,启动部署任务")
				task.AddNewTask(string(bodyContent))
			} else {
				utils.Log2file("验证失败")
				_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
			}
		} else {
			_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Unmatch x-github-event\"}")
		}
	} else {
		_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method or unknown request url\"}")
	}
}
