package ai

import (
	"github.com/sashabaranov/go-openai"
)

type Openai struct {
	*openai.Client
	model string
}

func NewOpenAi(token, baseUrl, model string) *Openai {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseUrl

	return &Openai{
		Client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

func (o *Openai) Engine() string {
	return EngineOpenai
}

func (o *Openai) Chat() (string, error) {
	//stream, err := c.Openai.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
	//	Model:    c.model,
	//	Messages: make([]openai.ChatCompletionMessage, 0),
	//	Stream:   true,
	//})
	//response, err := stream.Recv()
	//for _, v := range response.Choices {
	//	//v.Delta.Content
	//}
	return "", nil
}
