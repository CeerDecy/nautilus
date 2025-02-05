package mq

const (
	ServerTypeMQTT = "mqtt"
)

type Config struct {
	Host       string `yaml:"host" env:"MQ_HOST"`
	Port       string `yaml:"port" env:"MQ_PORT"`
	Topic      string `yaml:"topic" env:"MQ_TOPIC"`
	ServerType string `yaml:"server_type" env:"MQ_SERVER_TYPE"`
}

type MessageHandler func(client Client, message []byte)

type Client interface {
	Subscribe(topic string, handler MessageHandler)
	Publish(topic string, data []byte) error
}
