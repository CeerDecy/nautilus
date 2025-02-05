package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"nautilus/nautilus-common/agent"
	"nautilus/nautilus-common/mq"
	"os"
	"reflect"
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
	if err := loadFromEnv(cfg.Agent); err != nil {
		return nil, fmt.Errorf("load Agent env failed: %v", err)
	}
	if err := loadFromEnv(cfg.MQ); err != nil {
		return nil, fmt.Errorf("load MQ env failed: %v", err)
	}
	return &cfg, nil
}

func loadFromEnv(obj any) error {
	of := reflect.TypeOf(obj)
	vf := reflect.ValueOf(obj)
	if of.Kind() == reflect.Ptr {
		of = of.Elem()
	}
	if vf.Kind() == reflect.Ptr {
		vf = vf.Elem()
	}
	for i := range of.NumField() {
		key := of.Field(i).Tag.Get("env")
		if key == "" {
			continue
		}
		if value := os.Getenv(key); value != "" {
			vf.Field(i).SetString(value)
		}
	}
	return nil
}
