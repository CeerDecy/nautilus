package ai

type DeepSeek struct {
}

func NewDeepSeek(token, baseUrl, model string) *DeepSeek {
	return &DeepSeek{}
}

func (d *DeepSeek) Engine() string {
	return EngineDeepSeek
}
