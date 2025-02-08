package agent

import (
	"context"
	"encoding/json"
	"github.com/erda-project/erda-infra/base/servicehub"
	"nautilus/nautilus-common/ai"
	"nautilus/nautilus-common/ai/agent/service"
	"nautilus/nautilus-common/ai/message"
	"nautilus/nautilus-common/mysql"
)

type config struct {
	Component string   `file:"Component"`
	Prompt    []string `file:"Prompt"`
}

type provider struct {
	Cfg           *config
	Ai            ai.Interface    `autowired:"nautilus-ai"`
	Db            mysql.Interface `autowired:"mysql-provider"`
	msg           chan string
	component     string
	conversations map[string]*message.Conversation
	resp          chan string
	prompt        []message.Message
	toolsService  *service.ToolsService
	tools         []ai.Tool
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	p.component = p.Cfg.Component
	p.toolsService = service.NewToolsService(p.Db.DB())
	tools := p.toolsService.GetToolsByRole(p.component)
	for _, tool := range tools {
		p.tools = append(p.tools, ai.Tool{
			Type: ai.ToolTypeFunction,
			Function: &ai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Strict:      tool.Strict,
				Parameters:  json.RawMessage(tool.Parameters),
			},
		})
	}
	p.Ai.SetTools(p.tools)
	for _, v := range p.Cfg.Prompt {
		p.prompt = append(p.prompt, message.Message{
			Role:    message.ChatMessageRoleSystem,
			Content: v,
		})
	}
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	p.Start()
	return nil
}

func init() {
	servicehub.Register("nautilus-ai-agent", &servicehub.Spec{
		Services:     []string{"nautilus-ai-agent"},
		Dependencies: []string{"nautilus-ai", "mysql-provider"},
		ConfigFunc:   func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{
				msg:           make(chan string),
				conversations: make(map[string]*message.Conversation),
				resp:          make(chan string),
				prompt:        make([]message.Message, 0),
				tools:         make([]ai.Tool, 0),
			}
		},
	})
}
