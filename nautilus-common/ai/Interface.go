package ai

import (
	"nautilus/nautilus-common/ai/message"
)

type Interface interface {
	Engine() string
	Send(conversation *message.Conversation) (message.Session, error)
	SetTools(tools []Tool)
}
