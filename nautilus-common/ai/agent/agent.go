package agent

import "nautilus/nautilus-common/ai"

type Client struct {
	AI  ai.Interface
	msg chan string
}

func NewAgent(ai ai.Interface) *Client {
	return &Client{
		AI:  ai,
		msg: make(chan string),
	}
}

func (c *Client) Send(msg string) {
	c.msg <- msg
}
