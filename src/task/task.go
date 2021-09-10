package task

import (
	"GoWebhook/src/utils"
	"fmt"
	"os/exec"
)

var running = false
var queue []*structTaskQueue

type structTaskQueue struct {
	script string
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
		script := queue[0].script
		queue = queue[:0:0]
		go startTask(script)
	}
}

func startTask(script string) {
	running = true
	cmd := exec.Command("/bin/sh", script)
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
