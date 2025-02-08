package captain

import (
	"context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"nautilus/nautilus-common/ai/agent"
	"nautilus/nautilus-common/k8s"
	"nautilus/nautilus-common/mq"
	"slices"
	"time"
)

type Core struct {
	mq               mq.Interface
	ai               agent.Interface
	k8s              k8s.Interface
	console          *Console
	ignoreNamespaces []string
	ignorePods       []string
}

func NewCore(m mq.Interface, a agent.Interface, k k8s.Interface) *Core {
	conResp := make(chan string)
	console := NewConsole(context.Background(), k)
	console.Start(conResp)
	return &Core{
		mq:               m,
		ai:               a,
		console:          console,
		k8s:              k,
		ignoreNamespaces: []string{"kube-system", "kubeprober"},
		ignorePods:       []string{},
	}
}

func (core *Core) Start() {
	go core.monitor()
}

func (core *Core) Stop() {}

func (core *Core) monitor() {
	tick := time.Tick(1 * time.Minute)
	for {
		if core.k8s == nil || core.k8s.ClientSet() == nil {
			logrus.Printf("k8s not configured ======> ")
		}
		list, err := core.k8s.ClientSet().CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			logrus.Printf("Failed to list pods: %v", err)
		}

		unHealthy := make([]corev1.Pod, 0)

		for _, pod := range list.Items {
			if slices.Index(core.ignoreNamespaces, pod.Namespace) >= 0 {
				continue
			}
			if pod.Status.Phase != corev1.PodRunning && pod.Status.Phase != corev1.PodSucceeded {
				logrus.Infof("%s Pod %s is in phase %s", pod.Namespace, pod.Name, pod.Status.Phase)
				unHealthy = append(unHealthy, pod)
			}
		}
		<-tick
	}
}
