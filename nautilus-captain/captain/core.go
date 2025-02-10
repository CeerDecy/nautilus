package captain

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"nautilus/nautilus-common/ai/agent"
	"nautilus/nautilus-common/ai/model"
	"nautilus/nautilus-common/k8s"
	"nautilus/nautilus-common/mq"
	"nautilus/nautilus-common/tools/markdown"
	"nautilus/nautilus-common/tools/str"
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

		unHealthy = make([]corev1.Pod, 0)

		var content string = "all pods are healthy"
		if len(unHealthy) != 0 {
			content = markdown.ToMarkdownTable(convertToPodMaps(unHealthy))
		}

		agentContent := model.AgentContent{
			Content: content,
			Id:      str.GenerateUUID(),
			Time:    time.Now().Format(time.RFC3339),
		}
		marshal, err := json.Marshal(agentContent)
		if err != nil {
			logrus.Printf("Failed to marshal content: %v", err)
			continue
		}
		core.ai.Send(string(marshal))
		//fmt.Println(table)
		<-tick
	}
}

func convertToPodMaps(pods []corev1.Pod) (res []map[string]string) {
	res = make([]map[string]string, 0, len(pods))
	for _, pod := range pods {
		podMap := make(map[string]string)
		podMap["namespace"] = pod.Namespace
		podMap["name"] = pod.Name
		podMap["phase"] = string(pod.Status.Phase)
		podMap["startTime"] = "none"
		if pod.Status.StartTime != nil {
			podMap["startTime"] = pod.Status.StartTime.Format(time.RFC3339)
		}
		bytes, err := json.Marshal(pod.Status.ContainerStatuses)
		if err != nil {
			logrus.Errorf("Failed to marshal container statuses: %v", err)
		}
		podMap["containerStatuses"] = string(bytes)
		res = append(res, podMap)
	}
	return res
}
