package captain

import (
	"context"
	"nautilus/nautilus-common/ai"
	"nautilus/nautilus-common/k8s"
	"nautilus/nautilus-common/mq"
)

type Core struct {
	Mq      mq.Interface
	Ai      ai.Interface
	console *Console
}

func NewCore(m mq.Interface, a ai.Interface, k k8s.Interface) *Core {
	conResp := make(chan string)
	console := NewConsole(context.Background(), k)
	console.Start(conResp)
	return &Core{
		Mq:      m,
		Ai:      a,
		console: console,
	}
}

func (core *Core) Start() {
	go core.monitor()
}

func (core *Core) Stop() {}

func (core *Core) monitor() {

}
