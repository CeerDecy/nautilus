package main

import (
	_ "embed"
	"github.com/erda-project/erda-infra/base/servicehub"

	_ "github/ceerdecy/nautilus/nautilus-captain/captain"
	_ "github/ceerdecy/nautilus/nautilus-common/ai"
	_ "github/ceerdecy/nautilus/nautilus-common/ai/agent"
	_ "github/ceerdecy/nautilus/nautilus-common/k8s"
	_ "github/ceerdecy/nautilus/nautilus-common/mq"
	_ "github/ceerdecy/nautilus/nautilus-common/mysql"
)

//go:embed bootstrap.yaml
var bootstrap string

func main() {
	hub := servicehub.New()
	hub.RunWithOptions(&servicehub.RunOptions{
		Content: bootstrap,
	})
}
