package main

import (
	_ "embed"
	"fmt"
	"nautilus/nautilus-captain/config"
	"nautilus/nautilus-captain/message"
	"nautilus/nautilus-common/mq"
	"time"

	"github.com/sirupsen/logrus"
)

//go:embed bootstrap.yaml
var bootstrap string

func main() {
	cfg, err := config.LoadConfig(bootstrap)
	if err != nil {
		logrus.Fatal(err)
	}
	client, err := message.NewClient(cfg.MQ)
	if err != nil {
		logrus.Fatal(err)
	}
	client.Subscribe("captain", func(client mq.Client, message []byte) {
		fmt.Println(string(message))
	})

	go func() {
		for {
			time.Sleep(5 * time.Second)
			_ = client.Publish("captain", []byte(`{"command": "list tasks", "args": []"}`))
		}
	}()

	channel := make(chan int)
	<-channel
}
