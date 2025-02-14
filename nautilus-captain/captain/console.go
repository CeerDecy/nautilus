package captain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github/ceerdecy/nautilus/nautilus-common/k8s"
	"sync"
)

const (
	CommandCreate = "create"
)

var tasks = &sync.Map{}

type Command struct {
	Name string
	Args []string
}

func (c *Command) doCreate() error {
	if len(c.Args) < 1 {
		return fmt.Errorf("no command specified")
	}
	tpy := c.Args[0]
	switch tpy {
	case "task":
		tasks.Store(c.Args[1], c.Args[2])
	}
	return nil
}

type Console struct {
	K8s k8s.Interface
	cmd chan Command
	ctx context.Context
}

func NewConsole(ctx context.Context, k k8s.Interface) *Console {
	return &Console{
		K8s: k,
		ctx: ctx,
	}
}

func (c *Console) Start(resp chan string) {
	go func() {
		select {
		case <-c.ctx.Done():
			return
		case cmd := <-c.cmd:
			dict := make(map[string]any)
			switch cmd.Name {
			case "create":
				err := cmd.doCreate()
				if err != nil {
					dict["err"] = err.Error()
				}
				res, err := json.Marshal(dict)
				if err != nil {
					logrus.Errorf(err.Error())
				}
				resp <- string(res)
			}
		}
	}()
}
