package dto

import "Calculator/internal/executor"

type GroupedInstructions struct {
	Expressions []executor.Expression
	VarsToPrint map[string]bool
}
