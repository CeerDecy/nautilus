package ai

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"nautilus/nautilus-common/ai/message"
	"sync"
)

type Openai struct {
	*openai.Client
	model string
	tools []openai.Tool
}

func NewOpenAi(token, baseUrl, model string) *Openai {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseUrl

	return &Openai{
		Client: openai.NewClientWithConfig(config),
		model:  model,
		tools:  make([]openai.Tool, 0),
	}
}

func (o *Openai) Engine() string {
	return EngineOpenai
}

func (o *Openai) SetTools(tools []Tool) {
	for _, v := range tools {
		o.tools = append(o.tools, openai.Tool{
			Type: openai.ToolType(v.Type),
			Function: &openai.FunctionDefinition{
				Name:        v.Function.Name,
				Description: v.Function.Description,
				Strict:      v.Function.Strict,
				Parameters:  v.Function.Parameters,
			},
		})
	}
}

func (o *Openai) Send(conversation *message.Conversation) (message.Session, error) {
	openaiMessage := convertToOpenaiMessage(conversation)
	fmt.Printf("function tools ===> %v", o.tools)
	stream, err := o.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: openaiMessage,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Stream: true,
		Tools:  o.tools,
	})
	if err != nil {
		return nil, err
	}
	return NewOpenAiSession(stream), nil
}

type OpenaiSession struct {
	stream  *openai.ChatCompletionStream
	buf     []byte
	content []byte
	done    chan struct{}
	isDone  bool
	lock    sync.Mutex
}

func NewOpenAiSession(stream *openai.ChatCompletionStream) *OpenaiSession {
	session := &OpenaiSession{
		stream:  stream,
		buf:     make([]byte, 0),
		content: make([]byte, 0),
		done:    make(chan struct{}, 1),
		isDone:  false,
	}
	go session.readStream()
	return session
}

func (o *OpenaiSession) readBuf() []byte {
	o.lock.Lock()
	defer o.lock.Unlock()
	if len(o.buf) == 0 {
		return []byte{}
	}
	res := make([]byte, len(o.buf))
	copy(res, o.buf)
	o.buf = make([]byte, 0)
	return res
}

func (o *OpenaiSession) readStream() {
	for {
		fmt.Printf("readStream start\n")
		recv, err := o.stream.Recv()
		if err != nil {
			fmt.Printf("readStream recv err: %v\n", err)
			o.lock.Lock()
			fmt.Printf("content ====> %v\n", string(o.content))
			_ = o.stream.Close()
			o.isDone = true
			o.done <- struct{}{}
			o.lock.Unlock()
		}
		//recv.ID
		for _, v := range recv.Choices {
			o.lock.Lock()
			o.buf = append(o.buf, v.Delta.Content...)
			o.content = append(o.content, v.Delta.Content...)
			for _, toolcall := range v.Delta.ToolCalls {
				//if toolcall.Function.Name != "" {
				//
				//}
				fmt.Printf("toolcall ===> %+v\n", toolcall)
			}
			o.lock.Unlock()
		}
	}
}

func (o *OpenaiSession) ReadMessage() []byte {
	if o.isDone {
		return o.content
	}
	write := o.HandleWrite()

	var buffer bytes.Buffer
	for write(&buffer) {
	}
	return buffer.Bytes()
}

func (o *OpenaiSession) HandleWrite() func(writer io.Writer) bool {
	return func(writer io.Writer) bool {
		o.lock.Lock()
		defer o.lock.Unlock()
		buf := o.readBuf()
		var err error
		if len(buf) > 0 {
			_, err = writer.Write(buf)
		}
		select {
		case <-o.done:
			return false
		default:
			return err == nil
		}
	}
}

func convertToOpenaiMessage(conversation *message.Conversation) []openai.ChatCompletionMessage {
	var openaiMessage []openai.ChatCompletionMessage
	for _, msg := range conversation.Messages() {
		openaiMessage = append(openaiMessage, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}
	return openaiMessage
}
