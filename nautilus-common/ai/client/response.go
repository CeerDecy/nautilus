package client

type RespType string

const RespMessage RespType = "message"
const RespToolCall RespType = "tool_call"

type Response struct {
	Type           RespType
	Content        []byte
	ConversationId string
	ToolCalls      []ToolCall
}
