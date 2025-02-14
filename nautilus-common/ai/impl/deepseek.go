package impl

import (
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
)

type DeepSeek struct {
}

func (d *DeepSeek) SetTools(tools []client.Tool) {
	//TODO implement me
	panic("implement me")
}

func (d *DeepSeek) Send(conversation *client.Conversation) (client.Session, error) {
	//TODO implement me
	panic("implement me")
}

func NewDeepSeek(token, baseUrl, model string) client.Interface {
	return &DeepSeek{}
}

func (d *DeepSeek) Engine() string {
	return client.EngineDeepSeek
}
