package task

import (
	"GoWebhook/src/config"
	"GoWebhook/src/utils"
	"fmt"
	"os/exec"
)

var running = false
var queue []*structTaskQueue

type structTaskQueue struct {
	requestBodyString string
}

// AddNewTask add new task
func AddNewTask(bodyContent string) {
	queue = append(queue, newStructTaskQueue(bodyContent))
	checkoutTaskStatus()
}

func newStructTaskQueue(body string) *structTaskQueue {
	return &structTaskQueue{body}
}

func checkoutTaskStatus() {
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
		checkoutTaskStatus()
	} else {
		running = false
		utils.Log2file(fmt.Sprintf("部署失败:\n %s", err))
		checkoutTaskStatus()
	}
}
