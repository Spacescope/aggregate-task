package tasks

import (
	"aggregate-task/internal/busi/core/tasks/evmtask"
	"context"
	"fmt"
)

var taskMap = make(map[string]Task)

type Task interface {
	Name() string
	Model() interface{}
	Run(ctx context.Context, height int64 /*, version int*/) error
}

func Register(tasks ...Task) {
	for _, task := range tasks {
		_, ok := taskMap[task.Name()]
		if ok {
			panic(fmt.Sprintf("task nae %s conflict", task.Name()))
		}
		taskMap[task.Name()] = task
	}
}

func GetTask(name string) Task {
	return taskMap[name]
}

func init() {
	Register(new(evmtask.GasOutput))
}
