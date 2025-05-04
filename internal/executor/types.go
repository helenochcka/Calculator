package executor

type ResultProcessor func(result Result) (bool, error)
