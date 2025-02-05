package message

import (
	"fmt"
	"nautilus/nautilus-common/mq"
	"nautilus/nautilus-common/mq/mqtt"
)

type Client struct {
	mq.Client
}

func NewClient(cfg *mq.Config) (*Client, error) {
	var client mq.Client
	switch cfg.ServerType {
	case mq.ServerTypeMQTT:
		client = mqtt.NewClient(cfg.Host, cfg.Port, "nautilus-captain")
	default:
		return nil, fmt.Errorf("not supported server type %s", cfg.ServerType)
	}
	return &Client{
		Client: client,
	}, nil
}
