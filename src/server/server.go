package server

import (
	"GoWebhook/src/config"
	"GoWebhook/src/task"
	"GoWebhook/src/utils"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// StartService start service
func StartService(address string, port string) error {

	router := mux.NewRouter()

	router.Methods(http.MethodOptions)

	router.HandleFunc("/", index).Methods(http.MethodGet)

	var paths []string

	for _, server := range config.Instance.Servers {
		if !contains(paths, server.RequestUrl) {
			router.HandleFunc(server.RequestUrl, autoBuild).Methods(http.MethodPost)
			paths = append(paths, server.RequestUrl)
			utils.Log2file(fmt.Sprintf("append: %s", server.RequestUrl))
		} else {
			utils.Log2file(fmt.Sprintf("skip: %s", server.RequestUrl))
		}
	}

	utils.Log2file(fmt.Sprintf("service starting... %s:%s", address, port))
	return http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router)
}

func index(w http.ResponseWriter, r *http.Request) {
	utils.Log2file(r.URL.Host)
	_, _ = fmt.Fprintln(w, "{\"code\":200, \"description\":\"service running...\"}")
}

func autoBuild(w http.ResponseWriter, r *http.Request) {

	var servers []config.Server
	for _, server := range config.Instance.Servers {
		if server.RequestUrl == r.URL.RequestURI() {
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Error Method or unknown request url\"}")
		return
	}

	matchCount := 0
	for _, server := range servers {
		if r.Header.Get("x-github-event") == "push" {
			bodyContent, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			signature := r.Header.Get("X-Hub-Signature")
			if utils.VerifySignature(signature, string(bodyContent), server.Secret) {
				_, _ = fmt.Fprintln(w, "{\"code\":200, \"description\":\"OK\"}")
				utils.Log2file("验证通过,启动部署任务")
				task.AddNewTask(server.Identifier, server.Script)
				break
			} else {
				utils.Log2file("验证失败")
				_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Signature error\"}")
			}
			matchCount++
		}
	}
	if matchCount == 0 {
		_, _ = fmt.Fprintln(w, "{\"code\":200, \"error\":\"Miss match x-github-event\"}")
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
