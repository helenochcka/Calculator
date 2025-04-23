package executor

type Instruction struct {
	Type      string  `json:"type"`
	Operation *string `json:"op"`
	Result    string  `json:"var"`
	Left      *string `json:"left"`
	Right     *string `json:"right"`
}

type Item struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}
