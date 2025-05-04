package services

import (
	"Calculator/internal/executor"
	"fmt"
)

type ValidationService struct {
}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (vs *ValidationService) ValidateArgType(arg interface{}) error {
	switch arg.(type) {
	case string:
		return nil
	case int:
		return nil
	default:
		return fmt.Errorf("%w: %v", executor.ErrUnsupportedArgType, arg)
	}
}

func (vs *ValidationService) CheckIfVarAlreadyUsed(expressionVars map[string]bool, variable string) error {
	if _, ok := expressionVars[variable]; ok {
		return fmt.Errorf("%w: %v", executor.ErrVarAlreadyUsed, variable)
	}
	return nil
}

func (vs *ValidationService) CheckCyclicDependency(expressions []executor.Expression, variable string) error {
	if expressions != nil {
		for _, expression := range expressions {
			if expression.Variable == variable {
				return executor.ErrCyclicDependency
			}
		}
	}
	return nil
}

func (vs *ValidationService) CheckIfArgNeverCalculated(
	dependencyMap map[string][]executor.Expression,
	expressionVars map[string]bool,
) error {
	for variable := range dependencyMap {
		if _, exists := expressionVars[variable]; !exists {
			return fmt.Errorf("%w: %v", executor.ErrVarWillNeverBeCalc, variable)
		}
	}
	return nil
}

func (vs *ValidationService) CheckIfPrintVarNeverCalculated(
	varsToPrint map[string]bool,
	expressionVars map[string]bool,
) error {
	for variable := range varsToPrint {
		if _, exists := expressionVars[variable]; !exists {
			return fmt.Errorf("%w: %v", executor.ErrVarWillNeverBeCalc, variable)
		}
	}
	return nil
}
