package captain

import (
	"context"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/sirupsen/logrus"
	"nautilus/nautilus-common/ai"
	"nautilus/nautilus-common/mq"
)

type cfg struct {
}
type provider struct {
	Cfg *cfg
	Mq  mq.Interface `autowired:"nautilus-mq"`
	Ai  ai.Interface `autowired:"nautilus-ai"`
}

func (p *provider) Init(ctx servicehub.Context) (err error) {
	logrus.Infof("nautilus captain provider init")
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	return nil
}

func init() {
	servicehub.Register("captain", &servicehub.Spec{
		Services: []string{"captain"},
		Dependencies: []string{
			"nautilus-mq",
			"nautilus-ai",
		},
		OptionalDependencies: []string{},
		ConfigFunc:           func() interface{} { return &cfg{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
