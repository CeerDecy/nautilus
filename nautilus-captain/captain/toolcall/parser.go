package toolcall

import (
	"fmt"
)

const (
	ListTasks  = "list_tasks"
	CreateTask = "create_task"
	None       = "none"
)

func ParseToolCall(name string, args string) (*Reply, error) {
	switch name {
	case ListTasks:
	case CreateTask:
	case None:
		return &Reply{}, nil
	default:
		return &Reply{
			Err: fmt.Sprintf("unknown tool call: %s", name),
		}, nil
	}
	return nil, nil
}

func list(tpy string, args string) {

}
