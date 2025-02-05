package agent

import (
	"fmt"
	"nautilus/nautilus-common/agent/deepseek"
	"nautilus/nautilus-common/agent/openai"
)

const (
	DeepSeek = "deepseek"
	Openai   = "openai"
)

type Config struct {
	Engine  string `yaml:"engine"`
	Token   string `yaml:"token"`
	BaseUrl string `yaml:"baseUrl"`
	Model   string `yaml:"model"`
}

type Agent interface{}

func NewAgent(config Config) (Agent, error) {
	switch config.Engine {
	case DeepSeek:
		return deepseek.NewClient(config.Token, config.BaseUrl, config.Model), nil
	case Openai:
		return openai.NewClient(config.Token, config.BaseUrl, config.Model), nil
	}
	return nil, fmt.Errorf("engine %s not support", config.Engine)
}
