package main

import (
	"GoWebhook/src/config"
	"GoWebhook/src/server"
	"log"
)

func main() {

	config.LoadConfig2()
	//if errorString := config.LoadConfig(); errorString != "" {
	//	utils.Log2file(errorString)
	//	os.Exit(1)
	//}
	//
	listenErr := server.StartService(config.Instance.Host, config.Instance.Port)
	if listenErr != nil {
		log.Fatal("ListenAndServer error: ", listenErr)
	} else {
		println("ListenAndServer")
	}

}
