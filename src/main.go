package main

import (
	"GoWebhooks/src/config"
	"GoWebhooks/src/server"
	"GoWebhooks/src/utils"
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
