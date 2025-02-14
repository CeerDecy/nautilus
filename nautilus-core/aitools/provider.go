package aitools

import (
	"github.com/erda-project/erda-infra/base/servicehub"
	"github/ceerdecy/nautilus/nautilus-common/mysql"
)

type config struct{}
type provider struct {
	Cfg   *config
	Mysql mysql.Interface `autowired:"mysql-provider"`
}

func (p *provider) Init(ctx servicehub.Context) error {
	return nil
}

func init() {
	servicehub.Register("ai-tools", &servicehub.Spec{
		Services:    []string{"ai-tools"},
		Description: "ai-tools service",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
