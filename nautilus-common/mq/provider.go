package mq

import (
	"fmt"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/sirupsen/logrus"
)

const (
	ServerTypeMQTT = "mqtt"
)

type Config struct {
	Host       string `yaml:"host" env:"MQ_HOST" file:"host"`
	Port       string `yaml:"port" env:"MQ_PORT" file:"port"`
	Topic      string `yaml:"topic" env:"MQ_TOPIC" file:"topic"`
	ServerType string `yaml:"server_type" env:"MQ_SERVER_TYPE" file:"server_type"`
}

type provider struct {
	Cfg *Config
	Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	switch p.Cfg.ServerType {
	case ServerTypeMQTT:
		p.Interface = NewMQTT(p.Cfg.Host, p.Cfg.Port, p.Cfg.Topic)
	default:
		return fmt.Errorf("not supported server type: %s", p.Cfg.ServerType)
	}
	logrus.Infof("mq provider initialized")
	return nil
}

func init() {
	servicehub.Register("nautilus-mq", &servicehub.Spec{
		Services:             []string{"nautilus-mq"},
		Dependencies:         []string{},
		OptionalDependencies: []string{},
		ConfigFunc:           func() interface{} { return &Config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
