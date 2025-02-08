package message

import "github.com/sirupsen/logrus"

type Conversation struct {
	prompt []Message
	msg    []Message
	maxLen int
}

func NewConversation(prompt []Message, maxLen int) *Conversation {
	if prompt == nil {
		prompt = make([]Message, 0)
	}
	if maxLen < 0 {
		maxLen = DefaultMaxMessageLen
	}
	c := &Conversation{prompt: prompt, msg: make([]Message, 0, maxLen), maxLen: maxLen}
	return c
}

func (c *Conversation) Append(role Role, content string) {
	msg := append(c.msg, Message{
		Role:    role,
		Content: content,
	})
	if len(msg) > c.maxLen {
		msg = msg[len(msg)-c.maxLen:]
	}
	c.msg = msg
}

func (c *Conversation) Messages() []Message {
	logrus.Infof("conversation:%v", len(c.msg))
	logrus.Infof("conversation:%v", len(c.prompt))
	return append(c.prompt, c.msg...)
}
