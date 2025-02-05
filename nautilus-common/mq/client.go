package mq

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Topic      string `yaml:"topic"`
	ServerType string `yaml:"server_type"`
}
