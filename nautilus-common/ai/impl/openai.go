package impl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github/ceerdecy/nautilus/nautilus-common/ai/client"
	"io"
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
	return client.EngineOpenai
}

func (o *Openai) SetTools(tools []client.Tool) {
	for _, v := range tools {
		if v.Function == nil {
			continue
		}
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

func (o *Openai) Send(conversation *client.Conversation) (client.Session, error) {
	openaiMessage := convertToOpenaiMessage(conversation)
	fmt.Printf("function tools ===> %v\n", o.tools)
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
	stream    *openai.ChatCompletionStream
	buf       []byte
	content   []byte
	done      chan struct{}
	isDone    bool
	lock      sync.Mutex
	role      string
	refusal   string
	toolCalls []openai.ToolCall
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
		recv, err := o.stream.Recv()
		if err != nil {
			fmt.Printf("stream.Recv() err: %v\n", err)
			o.lock.Lock()
			_ = o.stream.Close()
			o.isDone = true
			o.done <- struct{}{}
			o.lock.Unlock()
			return
		}
		//recv.ID
		for _, v := range recv.Choices {
			o.lock.Lock()
			o.buf = append(o.buf, v.Delta.Content...)
			o.content = append(o.content, v.Delta.Content...)
			for _, toolcall := range v.Delta.ToolCalls {
				o.toolCalls = append(o.toolCalls, toolcall)
			}
			o.role = v.Delta.Role
			o.refusal = v.Delta.Refusal
			o.lock.Unlock()
		}
	}
}

func (o *OpenaiSession) convertToolCalls() []client.ToolCall {
	var toolCalls = make([]client.ToolCall, 0)
	for _, call := range o.toolCalls {
		toolCalls = append(toolCalls, client.ToolCall{
			Type: client.ToolType(call.Type),
			Function: client.FunctionCall{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
			ID:    call.ID,
			Index: call.Index,
		})
	}
	return toolCalls
}

func (o *OpenaiSession) ReadMessage() client.Message {
	if o.isDone {
		return client.Message{
			Role:      client.Role(o.role),
			Content:   o.content,
			Refusal:   o.refusal,
			ToolCalls: o.convertToolCalls(),
		}
	}
	write := o.HandleWrite()

	var buffer bytes.Buffer
	for write(&buffer) {
	}
	return client.Message{
		Role:      client.Role(o.role),
		Content:   buffer.Bytes(),
		Refusal:   o.refusal,
		ToolCalls: o.convertToolCalls(),
	}
}

func (o *OpenaiSession) HandleWrite() func(writer io.Writer) bool {
	return func(writer io.Writer) bool {
		buf := o.readBuf()
		var err error
		if len(buf) > 0 {
			_, err = writer.Write(buf)
			//fmt.Printf("buffer ====> %v\n", string(buf))
		}
		select {
		case <-o.done:
			return false
		default:
			return err == nil
		}
	}
}

func convertToOpenaiMessage(conversation *client.Conversation) []openai.ChatCompletionMessage {
	var openaiMessage []openai.ChatCompletionMessage
	for _, msg := range conversation.Messages() {
		openaiMessage = append(openaiMessage, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: string(msg.Content),
			Refusal: msg.Refusal,
			//ToolCalls: make([]openai.ToolCall, 0),
		})
	}
	return openaiMessage
}
