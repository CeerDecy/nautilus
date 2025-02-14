package client

const (
	EngineDeepSeek = "deepseek"
	EngineOpenai   = "openai"
)

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type ToolCall struct {
	Type     ToolType     `json:"type"`
	ID       string       `json:"id"`
	Function FunctionCall `json:"function"`
	Index    *int         `json:"index"`
}

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type Tool struct {
	Type     ToolType            `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
	Parameters  any    `json:"parameters"`
}
