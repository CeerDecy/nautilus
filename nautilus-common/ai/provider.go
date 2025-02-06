package ai

import (
	"fmt"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type provider struct {
	Cfg *Config
	Interface
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	switch p.Cfg.Engine {
	case EngineOpenai:
		p.Interface = NewOpenAi(p.Cfg.Token, p.Cfg.BaseUrl, p.Cfg.Model)
	case EngineDeepSeek:
		p.Interface = NewDeepSeek(p.Cfg.Token, p.Cfg.BaseUrl, p.Cfg.Model)
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
