package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
	"github/ceerdecy/nautilus/nautilus-common/ai/model"
)

type Interface interface {
	Send(data []byte) (chan client.Response, error)
}

func (p *provider) Send(data []byte) (chan client.Response, error) {
	var content model.AgentContent
	err := json.Unmarshal(data, &content)
	if err != nil {
		logrus.Errorf("unmarshal agent content error: %s", err.Error())
		return nil, err
	}
	p.response[content.Id] = make(chan client.Response)
	p.msg <- content
	return p.response[content.Id], nil
}

func (p *provider) Start() {
	for {
		content := <-p.msg
		conversation, ok := p.conversations[content.Id]
		if !ok {
			conversation = client.NewConversation(p.prompt, 10)
			p.conversations[content.Id] = conversation
		}
		conversation.Append(client.ChatMessageRoleUser, content.Content)
		session, err := p.Ai.Send(conversation)
		if err != nil {
			logrus.Errorf("send message to agent error: %s", err.Error())
		}
		_ = session

		message := session.ReadMessage()

		var respType = client.RespMessage
		if len(message.ToolCalls) > 0 {
			respType = client.RespToolCall
		}

		p.response[content.Id] <- client.Response{
			Content:        message.Content,
			ConversationId: content.Id,
			Type:           respType,
			ToolCalls:      message.ToolCalls,
		}
	}
}
