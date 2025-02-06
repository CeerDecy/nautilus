package mq

type MessageHandler func(client Interface, message []byte)

type Interface interface {
	Subscribe(topic string, handler MessageHandler)
	Publish(topic string, data []byte) error
}
