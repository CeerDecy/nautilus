package openai

import (
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	*openai.Client
	model string
}

func NewClient(token, baseUrl, model string) *Client {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseUrl

	return &Client{
		Client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

func (c *Client) Chat() (string, error) {
	//stream, err := c.Client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
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
