package task

import (
	"os/exec"
	"GoWebhooks/src/utils"
	"fmt"
	"GoWebhooks/src/config"
)

var running = false
var queue []*structTaskQueue

type structTaskQueue struct {
	requestBodyString string
}

func AddNewTask(bodyContent string)  {
	queue = append(queue, NewStructTaskQueue(bodyContent))
}

func NewStructTaskQueue(body string) *structTaskQueue {
	return &structTaskQueue{body}
}

func CheckoutTaskStatus()  {
	if running {
		return
	}
	if len(queue) > 0 {
		queue = queue[:0:0]
		go startTask()
	}
}

func startTask() {
	running = true
	cmd := exec.Command("/bin/sh", config.GetShell())
	_, err := cmd.Output()
	if err == nil {
		running = false
		utils.Log2file("部署成功")
		CheckoutTaskStatus()
	} else {
		running = false
		utils.Log2file(fmt.Sprintf("部署失败:\n %s", err))
		CheckoutTaskStatus()
	}
}