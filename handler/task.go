package handler

import (
	"sync"
	"time"
	"github.com/astaxie/beego/toolbox"
)

var mu = new(sync.RWMutex)

func AddTask(tname string, spec string, f toolbox.TaskFunc) {
	t := toolbox.NewTask(tname, spec, f)
	mu.Lock()
	toolbox.AddTask(tname, t)
	SetNext(t)
	mu.Unlock()
}

func RunTask() {
	toolbox.StartTask()
}

func SetNext(t *toolbox.Task) {
	t.SetNext(time.Now().Local())
}

func StopSpeciTask(tname string) {
	mu.Lock()
	toolbox.DeleteTask(tname)
	mu.Unlock()
}

func StopAll() {
	toolbox.StopTask()
}