package main

import (
	_ "embed"
	"github.com/erda-project/erda-infra/base/servicehub"

	_ "nautilus/nautilus-captain/captain"
	_ "nautilus/nautilus-common/ai"
	_ "nautilus/nautilus-common/k8s"
	_ "nautilus/nautilus-common/mq"
)

//go:embed bootstrap.yaml
var bootstrap string

func main() {
	hub := servicehub.New()
	hub.RunWithOptions(&servicehub.RunOptions{
		Content: bootstrap,
	})
}
