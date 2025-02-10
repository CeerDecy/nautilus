package agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"nautilus/nautilus-common/ai/message"
	"nautilus/nautilus-common/ai/model"
)

type Interface interface {
	Send(data string)
}

func (p *provider) Send(data string) {
	p.msg <- data
}

func (p *provider) Start() {
	for {
		msg := <-p.msg
		var content model.AgentContent
		err := json.Unmarshal([]byte(msg), &content)
		if err != nil {
			logrus.Errorf("unmarshal agent content error: %s", err.Error())
			continue
		}
		conversation, ok := p.conversations[content.Id]
		if !ok {
			conversation = message.NewConversation(p.prompt, 10)
			p.conversations[content.Id] = conversation
		}
		conversation.Append(message.ChatMessageRoleUser, content.Content)
		session, err := p.Ai.Send(conversation)
		if err != nil {
			logrus.Errorf("send message to agent error: %s", err.Error())
		}
		_ = session
		logrus.Infof("%+v", conversation.Messages())
	}
}
