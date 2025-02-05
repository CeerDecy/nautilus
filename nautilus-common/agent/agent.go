package agent

const (
	DeepSeek = "deepseek"
	Openai   = "openai"
)

type Config struct {
	Engine  string `yaml:"engine" env:"AGENT_ENGINE"`
	Token   string `yaml:"token" env:"AGENT_TOKEN"`
	BaseUrl string `yaml:"baseUrl" env:"AGENT_BASEURL"`
	Model   string `yaml:"model" env:"AGENT_MODEL"`
}

type Agent interface{}
