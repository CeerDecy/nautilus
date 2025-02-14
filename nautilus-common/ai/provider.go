package ai

import (
	"fmt"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
	"github/ceerdecy/nautilus/nautilus-common/ai/impl"
)

type Config struct {
	Engine  string `yaml:"engine" env:"AGENT_ENGINE" file:"engine"`
	Token   string `yaml:"token" env:"AGENT_TOKEN" file:"token"`
	BaseUrl string `yaml:"baseUrl" env:"AGENT_BASEURL" file:"baseUrl"`
	Model   string `yaml:"model" env:"AGENT_MODEL" file:"model"`
}

type provider struct {
	Cfg *Config
	client.Interface
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	switch p.Cfg.Engine {
	case client.EngineOpenai:
		p.Interface = impl.NewOpenAi(p.Cfg.Token, p.Cfg.BaseUrl, p.Cfg.Model)
	case client.EngineDeepSeek:
		p.Interface = impl.NewDeepSeek(p.Cfg.Token, p.Cfg.BaseUrl, p.Cfg.Model)
	default:
		return fmt.Errorf("ai engine %s not support", p.Cfg.Model)
	}
	return nil
}

func init() {
	servicehub.Register("nautilus-ai", &servicehub.Spec{
		Services:             []string{"nautilus-ai"},
		Dependencies:         []string{},
		OptionalDependencies: []string{},
		ConfigFunc:           func() interface{} { return &Config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
