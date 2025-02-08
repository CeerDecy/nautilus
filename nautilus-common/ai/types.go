package ai

const (
	EngineDeepSeek = "deepseek"
	EngineOpenai   = "openai"
)

type Config struct {
	Engine  string `yaml:"engine" env:"AGENT_ENGINE" file:"engine"`
	Token   string `yaml:"token" env:"AGENT_TOKEN" file:"token"`
	BaseUrl string `yaml:"baseUrl" env:"AGENT_BASEURL" file:"baseUrl"`
	Model   string `yaml:"model" env:"AGENT_MODEL" file:"model"`
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

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
