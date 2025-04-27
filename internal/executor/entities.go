package executor

type Instruction struct {
	Type      string  `json:"type"`
	Operation *string `json:"op"`
	Variable  string  `json:"var"`
	Left      *string `json:"left"`
	Right     *string `json:"right"`
}

type Item struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}

type Result struct {
	Lit    string `json:"lit"`
	Result int    `json:"result"`
}

type InstructionGin struct {
	Type      string      `json:"type"`
	Operation *string     `json:"op"`
	Variable  string      `json:"var"`
	Left      interface{} `json:"left"`
	Right     interface{} `json:"right"`
}

type Expression struct {
	Type      string
	Operation string
	Variable  string
	Left      interface{}
	Right     interface{}
}

type CalculationData struct {
	Operation string
	Variable  string
	Left      int
	Right     int
}
