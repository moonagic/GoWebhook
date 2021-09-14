package task

import (
	"GoWebhook/src/utils"
	"fmt"
	"os/exec"
)

var running = false

//var queue []*structTaskQueue

var queue = make(map[string][]*structTaskQueue)

type structTaskQueue struct {
	script string
}

// AddNewTask add new task
func AddNewTask(identifier string, script string) {

	if _, ok := queue[identifier]; ok {
		queue[identifier] = append(queue[identifier], newStructTaskQueue(script))
	} else {
		queue[identifier] = []*structTaskQueue{newStructTaskQueue(script)}
	}
	checkoutTaskStatus(identifier)
}

func newStructTaskQueue(body string) *structTaskQueue {
	return &structTaskQueue{body}
}

func checkoutTaskStatus(identifier string) {
	if running {
		return
	}
	if len(queue[identifier]) > 0 {
		script := queue[identifier][0].script
		queue[identifier] = queue[identifier][:0:0]
		go startTask(identifier, script)
	}
}

func startTask(identifier string, script string) {
	running = true
	cmd := exec.Command("/bin/sh", script)
	_, err := cmd.Output()
	if err == nil {
		running = false
		utils.Log2file("部署成功")
		checkoutTaskStatus(identifier)
	} else {
		running = false
		utils.Log2file(fmt.Sprintf("部署失败:\n %s", err))
		checkoutTaskStatus(identifier)
	}
}
