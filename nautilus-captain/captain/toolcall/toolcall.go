package toolcall

type Parser struct {
	Tools []Item
}

type Item struct {
	Name      string
	Parameter Parameter
}

type Reply struct {
	Msg string
	Err string
}

type Properties struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type Parameter struct {
	Type       string                  `json:"type"`
	Properties []map[string]Properties `json:"properties"`
	Required   []string                `json:"required"`
}
