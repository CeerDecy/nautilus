package captain

import (
	"context"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
)

type Agent struct {
	msg    chan string
	ai     client.Interface
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAgent(ctx context.Context, a client.Interface) *Agent {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &Agent{
		msg:    make(chan string),
		ai:     a,
		ctx:    ctx,
		cancel: cancelFunc,
	}
}

func (a *Agent) MessageChannel() chan string {
	return a.msg
}

func (a *Agent) Start() {
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				return
			case msg := <-a.msg:
				_ = msg
			}
		}
	}()
}

func (a *Agent) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
}
