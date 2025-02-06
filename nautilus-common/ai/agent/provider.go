package agent

import (
	"github.com/erda-project/erda-infra/base/servicehub"
	"nautilus/nautilus-common/ai"
	"nautilus/nautilus-common/ai/message"
)

type provider struct {
	ai            ai.Interface `autowire:"nautilus-ai"`
	msg           chan string
	conversations map[string]*message.Conversation
}

func (p *provider) Init() error {
	return nil
}

func (p *provider) Run() error {
	go p.Start()
	return nil
}

func init() {
	servicehub.Register("nautilus-ai-agent", &servicehub.Spec{
		Services:     []string{"nautilus-ai-agent"},
		Dependencies: []string{"nautilus-ai"},
		Description:  "nautilus-ai-agent",
		Creator: func() servicehub.Provider {
			return &provider{
				msg:           make(chan string),
				conversations: make(map[string]message.Conversation),
			}
		},
	})
}
