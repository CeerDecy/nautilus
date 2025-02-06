package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"nautilus/nautilus-common/ai/message"
	"nautilus/nautilus-common/ai/model"
)

func (p *provider) MessageChannel() chan string {
	return p.msg
}

func (p *provider) Start() {
	for {
		msg := <-p.msg
		var content model.AgentContent
		err := json.Unmarshal([]byte(msg), &content)
		if err != nil {
			logrus.Errorf("unmarshal agent content error: %s", err.Error())
		}
		conversation, ok := p.conversations[content.Id]
		if !ok {
			conversation = message.NewConversation([]message.Message{}, 10)
			p.conversations[content.Id] = conversation
		}
		conversation.Append(message.ChatMessageRoleUser, content.Content)
		p.ai.Send(conversation)
	}
}
