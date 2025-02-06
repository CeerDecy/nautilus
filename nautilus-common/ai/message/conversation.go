package message

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
	return &Conversation{prompt: prompt, msg: make([]Message, 0, maxLen), maxLen: maxLen}
}

func (c *Conversation) Append(role Role, content string) {
	msg := append(c.msg, Message{
		Role:    string(role),
		Content: content,
	})
	if len(msg) > c.maxLen {
		c.msg = msg[len(msg)-c.maxLen:]
	}
}

func (c *Conversation) Messages() []Message {
	return append(c.prompt, c.msg...)
}
