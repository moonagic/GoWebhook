package main

import (
	"GoWebhook/src/config"
	"GoWebhook/src/server"
	"GoWebhook/src/utils"
	"log"
	"os"
)

func main() {

	if errorString := config.LoadConfig(); errorString != "" {
		utils.Log2file(errorString)
		os.Exit(0)
	}

	listenErr := server.StartService(config.GetHost(), config.GetPort())
	if listenErr != nil {
		log.Fatal("ListenAndServer error: ", listenErr)
	}

}
