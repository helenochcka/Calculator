package executor

type Expression struct {
	Type      string
	Operation string
	Variable  string
	Left      interface{}
	Right     interface{}
}

type Result struct {
	Key   string
	Value int
}
