package executor

type MessageHandler func(result Result) (bool, error)
