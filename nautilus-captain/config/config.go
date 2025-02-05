package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"nautilus/nautilus-common/agent"
	"nautilus/nautilus-common/mq"
)

type Config struct {
	Agent *agent.Config `yaml:"agent"`
	MQ    *mq.Config    `yaml:"mq"`
}

func LoadConfig(bootstrap string) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal([]byte(bootstrap), &cfg); err != nil {
		return nil, fmt.Errorf("load config failed: %v", err)
	}
	return &cfg, nil
}
