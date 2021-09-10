package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var targetLogFile = "log"

// Log2file write log to file
func Log2file(content string) {

	targetLogDir := ""

	switch runtime.GOOS {
	case "linux":
		targetLogDir = "/etc/gowebhook/"
		break
	case "darwin":
		targetLogDir = "dev/"
		break
	}

	var err error

	if _, err := os.Stat(targetLogDir); err == nil {

	} else {
		fmt.Println("Dir not exists, try to create...", targetLogDir)
		err := os.MkdirAll(targetLogDir, 0711)
		if err != nil {
			fmt.Println("Error creating directory", targetLogDir)
			fmt.Println("err:", err)
			return
		}
	}

	if _, err := os.Stat(targetLogDir + targetLogFile); err == nil {

	} else {
		fmt.Println("Path not exists, try to create...", targetLogDir+targetLogFile)
		_, err := os.Create(targetLogDir + targetLogFile)
		if err != nil {
			fmt.Println("Error creating file", targetLogDir+targetLogFile)
			fmt.Println("err:", err)
			return
		}
	}

	f, err := os.OpenFile(targetLogDir+targetLogFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err == nil {
		timeString := time.Now().Format("2006-01-02 15:04:05")
		_, _ = f.WriteString("[" + timeString + "]" + "" + content)
		_, _ = f.WriteString("\n")
	} else {
		fmt.Println("Open file failed...", targetLogDir+targetLogFile)
		fmt.Println("err:", err)
	}
}
