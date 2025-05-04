package dto

import "Calculator/internal/executor"

type GroupedInstructions struct {
	Expressions []executor.Expression
	VarsToPrint map[string]bool
}

type CalculationData struct {
	Variable  string
	Operation string
	Left      int
	Right     int
	QueueName string
}
