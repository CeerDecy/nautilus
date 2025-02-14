package agent

import (
	"context"
	"encoding/json"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github/ceerdecy/nautilus/nautilus-common/ai/agent/service"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
	"github/ceerdecy/nautilus/nautilus-common/ai/model"
	"github/ceerdecy/nautilus/nautilus-common/mysql"
)

type config struct {
	Component string   `file:"Component"`
	Prompt    []string `file:"Prompt"`
}

type provider struct {
	Cfg           *config
	Ai            client.Interface `autowired:"nautilus-ai"`
	Db            mysql.Interface  `autowired:"mysql-provider"`
	msg           chan model.AgentContent
	component     string
	conversations map[string]*client.Conversation
	response      map[string]chan client.Response
	resp          chan string
	prompt        []client.Message
	toolsService  *service.ToolsService
	tools         []client.Tool
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	p.component = p.Cfg.Component
	p.toolsService = service.NewToolsService(p.Db.DB())
	tools := p.toolsService.GetToolsByRole(p.component)
	for _, tool := range tools {
		p.tools = append(p.tools, client.Tool{
			Type: client.ToolTypeFunction,
			Function: &client.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Strict:      tool.Strict,
				Parameters:  json.RawMessage(tool.Parameters),
			},
		})
	}
	p.Ai.SetTools(p.tools)
	for _, v := range p.Cfg.Prompt {
		p.prompt = append(p.prompt, client.Message{
			Role:    client.ChatMessageRoleSystem,
			Content: []byte(v),
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
				msg:           make(chan model.AgentContent),
				conversations: make(map[string]*client.Conversation),
				resp:          make(chan string),
				prompt:        make([]client.Message, 0),
				tools:         make([]client.Tool, 0),
				response:      make(map[string]chan client.Response),
			}
		},
	})
}
