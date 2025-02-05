package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"nautilus/nautilus-common/mq"
	"time"
)

type Client struct {
	mqtt.Client
	subTopics map[string]mq.MessageHandler
}

func NewClient(host, port, clientId string) *Client {
	var c = &Client{
		subTopics: make(map[string]mq.MessageHandler),
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	options.SetClientID(clientId)
	options.OnConnect = func(client mqtt.Client) {
		logrus.Infof("mqtt client connected, clientId=%s", clientId)
		for topic, handler := range c.subTopics {
			c.Subscribe(topic, handler)
		}
	}

	options.OnConnectionLost = func(client mqtt.Client, err error) {
		logrus.Errorf("mqtt client connection lost, clientId=%s, %s", clientId, err.Error())
	}

	options.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		logrus.Infof("mqtt client reconnecting, clientId=%s", clientId)
	}

	client := mqtt.NewClient(options)
	c.Client = client
	c.connect()
	return c
}

func (c *Client) Publish(topic string, data []byte) error {
	c.Client.Publish(topic, 2, false, data)
	return nil
}

func (c *Client) Subscribe(topic string, handler mq.MessageHandler) {
	if _, ok := c.subTopics[topic]; !ok {
		c.Unsubscribe(topic)
	}
	c.Client.Subscribe(topic, 2, func(client mqtt.Client, message mqtt.Message) {
		handler(c, message.Payload())
	})
	c.subTopics[topic] = handler
	logrus.Infof("mqtt client subscribed, topic=%s", topic)
}

func (c *Client) connect() {
	c.Connect()
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		if c.IsConnected() {
			return
		}
	}
}
