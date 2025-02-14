package message

import (
	"github.com/sirupsen/logrus"
	"github/ceerdecy/nautilus/nautilus-common/mq"
)

type Producer struct {
	client mq.Interface
	topic  string
	C      chan string
}

func NewProducer(topic string, mq mq.Interface) *Producer {
	return &Producer{
		client: mq,
		topic:  topic,
		C:      make(chan string),
	}
}

func (p *Producer) Start() {
	go func() {
		for {
			msg := <-p.C
			err := p.client.Publish(p.topic, []byte(msg))
			if err != nil {
				logrus.Errorf("publish fail %s", err)
			}
		}
	}()
}
