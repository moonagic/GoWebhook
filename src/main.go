package main

import (
	"GoWebhook/src/config"
	"GoWebhook/src/server"
	"github.com/fatih/color"
)

func main() {

	config.LoadConfig()

	color.Green("Service starting in %s:%s", config.Instance.Host, config.Instance.Port)
	listenErr := server.StartService(config.Instance.Host, config.Instance.Port)
	if listenErr != nil {
		color.Red("ListenAndServer error: ", listenErr)
	}

}
