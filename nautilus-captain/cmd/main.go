package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"nautilus/nautilus-captain/config"

	"github.com/sirupsen/logrus"
)

//go:embed bootstrap.yaml
var bootstrap string

func main() {
	cfg, err := config.LoadConfig(bootstrap)
	if err != nil {
		logrus.Fatal(err)
	}
	marshal, _ := json.Marshal(cfg)
	fmt.Println(string(marshal))
}
