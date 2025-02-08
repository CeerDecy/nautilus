package ai

import (
	"nautilus/nautilus-common/ai/message"
)

type DeepSeek struct {
}

func (d *DeepSeek) SetTools(tools []Tool) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) Send(conversation *message.Conversation) (message.Session, error) {
	//TODO implement me
	panic("implement me")
}

func NewDeepSeek(token, baseUrl, model string) Interface {
	return &DeepSeek{}
}

func (d *DeepSeek) Engine() string {
	return EngineDeepSeek
}
