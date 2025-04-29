package executor

type Instruction struct {
	Type      string      `json:"type"`
	Operation *string     `json:"op"`
	Variable  string      `json:"var"`
	Left      interface{} `json:"left"`
	Right     interface{} `json:"right"`
}

type Item struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}
