package executor

type Expression struct {
	Type      string
	Operation string
	Variable  string
	Left      any
	Right     any
}

type Result struct {
	Key   string
	Value int
}
