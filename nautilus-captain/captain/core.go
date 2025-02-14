package captain

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github/ceerdecy/nautilus/nautilus-captain/captain/toolcall"
	"github/ceerdecy/nautilus/nautilus-common/ai/agent"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
	"github/ceerdecy/nautilus/nautilus-common/ai/model"
	"github/ceerdecy/nautilus/nautilus-common/k8s"
	"github/ceerdecy/nautilus/nautilus-common/mq"
	"github/ceerdecy/nautilus/nautilus-common/tools/markdown"
	"github/ceerdecy/nautilus/nautilus-common/tools/str"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	parser           *toolcall.Parser
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

		//unHealthy = make([]corev1.Pod, 0)

		var content = "all pods are healthy"
		//var content string = "who are you?"
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
		responses, err := core.ai.Send(marshal)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
		_ = cancel
		go func(ctx context.Context) {
			select {
			case <-ctx.Done():
				logrus.Printf("Timeout waiting for response")
				return
			case resp := <-responses:
				core.dealResponse(resp)
			}
		}(ctx)
		<-tick
	}
}

func (core *Core) dealResponse(resp client.Response) {
	switch resp.Type {
	case client.RespMessage:
		logrus.Printf("conversation id [%v], response: %s", resp.ConversationId, string(resp.Content))
	case client.RespToolCall:
		for _, toolcall := range resp.ToolCalls {
			//aitools.RegisterToolCall(toolcall.Function.Name)
			toolcall.ParseToolCall(toolcall.Function.Name, toolcall.Function.Arguments)
		}
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
