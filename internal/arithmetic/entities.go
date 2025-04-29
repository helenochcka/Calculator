package arithmetic

type Expression struct {
	Op       string
	Variable string
	Left     int64
	Right    int64
}

type Result struct {
	Key   string
	Value int64
}
