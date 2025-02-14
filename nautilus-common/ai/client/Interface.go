package client

type Interface interface {
	Engine() string
	Send(conversation *Conversation) (Session, error)
	SetTools(tools []Tool)
}
