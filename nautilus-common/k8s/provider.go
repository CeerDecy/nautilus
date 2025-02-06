package k8s

import (
	"fmt"
	"github.com/erda-project/erda-infra/base/servicehub"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Interface interface {
	ClientSet() *kubernetes.Clientset
}

type config struct {
}
type provider struct {
	Cfg *config
	*kubernetes.Clientset
}

func (p *provider) Init(ctx servicehub.Context) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/Users/ceerdecy/.kube/config")
	if err != nil {
		return fmt.Errorf("fail to build config from flags: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("fail to build kubernetes client: %v", err)
	}
	p.Clientset = kubeClient
	return nil
}

func (p *provider) ClientSet() *kubernetes.Clientset {
	return p.Clientset
}

func init() {
	servicehub.Register("nautilus-kubernetes", &servicehub.Spec{
		Services:             []string{"nautilus-kubernetes"},
		Dependencies:         []string{},
		OptionalDependencies: []string{},
		ConfigFunc:           func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
