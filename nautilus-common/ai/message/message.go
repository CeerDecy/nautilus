package message

type Role string

const (
	ChatMessageRoleSystem    Role = "system"
	ChatMessageRoleUser      Role = "user"
	ChatMessageRoleAssistant Role = "assistant"

	DefaultMaxMessageLen = 10
)

type Message struct {
	Role    string
	Content string
}
