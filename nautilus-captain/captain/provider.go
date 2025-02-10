package captain

import (
	"context"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/sirupsen/logrus"
	"nautilus/nautilus-common/ai/agent"
	"nautilus/nautilus-common/k8s"
	"nautilus/nautilus-common/mq"
)

type config struct {
}
type provider struct {
	Cfg   *config
	Mq    mq.Interface    `autowired:"nautilus-mq"`
	Agent agent.Interface `autowired:"nautilus-ai-agent"`
	K8s   k8s.Interface   `autowired:"nautilus-kubernetes"`
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	logrus.Infof("nautilus captain provider init")
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	core := NewCore(p.Mq, p.Agent, p.K8s)
	core.Start()
	return nil
}

func init() {
	servicehub.Register("captain", &servicehub.Spec{
		Services: []string{"captain"},
		Dependencies: []string{
			"nautilus-mq",
			"nautilus-ai-agent",
			"nautilus-kubernetes",
		},
		OptionalDependencies: []string{},
		ConfigFunc:           func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
