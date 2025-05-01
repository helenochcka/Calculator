package gin_handlers

type Instruction struct {
	Type      string      `json:"type" binding:"required"`
	Operation *string     `json:"op"`
	Variable  string      `json:"var" binding:"required"`
	Left      interface{} `json:"left"`
	Right     interface{} `json:"right"`
}

type Item struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}
