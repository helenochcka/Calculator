package use_cases

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/services"
	"context"
	"fmt"
)

type UseCase struct {
	service *services.CommService
}

func NewUseCase(s *services.CommService) *UseCase {
	return &UseCase{service: s}
}

func (uc *UseCase) Execute(ctx context.Context, expressions *[]executor.Expression, varsToPrint map[string]bool) ([]executor.Result, error) {
	expressionMap := make(map[string][]executor.Expression)
	resultMap := make(map[string]int)
	expressionVars := make(map[string]bool)

	var prints []executor.Result

	reqId := ctx.Value("request_id").(string)
	if err := uc.service.DeclareQueue(reqId); err != nil {
		return nil, err
	}

	for _, expression := range *expressions {
		if _, ok := varsToPrint[expression.Variable]; ok {
			varsToPrint[expression.Variable] = true
		}
		if _, ok := expressionVars[expression.Variable]; ok {
			return nil, fmt.Errorf("%w: %v", executor.ErrVarAlreadyUsed, expression.Variable)
		}
		expressionVars[expression.Variable] = true
		left, leftIsInt := expression.Left.(int)
		right, rightIsInt := expression.Right.(int)
		if leftIsInt && rightIsInt {
			uc.service.RequestCalculation(reqId, left, right, expression.Variable, expression.Operation)
			continue
		}

		if !leftIsInt {
			if depExpressions := expressionMap[expression.Variable]; depExpressions != nil {
				for _, dep := range depExpressions {
					if dep.Variable == expression.Left.(string) {
						return nil, fmt.Errorf("%wfor variables: %v", executor.ErrCyclicDependency, fmt.Sprintf("%s and %s", expression.Variable, dep.Variable))
					}
				}
			}
			expressionMap[expression.Left.(string)] = append(expressionMap[expression.Left.(string)], expression)
		}
		if !rightIsInt {
			if depExpressions := expressionMap[expression.Variable]; depExpressions != nil {
				for _, dep := range depExpressions {
					if dep.Variable == expression.Right.(string) {
						return nil, fmt.Errorf("%wfor variables: %v", executor.ErrCyclicDependency, fmt.Sprintf("%s and %s", expression.Variable, dep.Variable))
					}
				}
			}
			expressionMap[expression.Right.(string)] = append(expressionMap[expression.Right.(string)], expression)
		}
	}

	for expressionVar := range expressionMap {
		if _, ok := expressionVars[expressionVar]; !ok {
			return nil, fmt.Errorf("%w: %v", executor.ErrVarNeverBeCalc, expressionVar)
		}
	}

	for variable, varToPrint := range varsToPrint {
		if !varToPrint {
			return nil, fmt.Errorf("%w: %v", executor.ErrVarToPrintNotFound, variable)
		}
	}

	handlerFunc := func(result executor.Result) (stop bool, err error) {
		resultMap[result.Key] = result.Value

		if varsToPrint[result.Key] {
			prints = append(prints, result)
			delete(varsToPrint, result.Key)
		}

		if len(varsToPrint) == 0 {
			return true, nil
		}

		if expressionMap[result.Key] != nil {
			for _, expression := range expressionMap[result.Key] {
				var left, right int

				switch expression.Left.(type) {
				case string:
					res, ok := resultMap[expression.Left.(string)]
					if !ok {
						continue
					}
					left = res
				case int:
					left = expression.Left.(int)
				}

				switch expression.Right.(type) {
				case string:
					res, ok := resultMap[expression.Right.(string)]
					if !ok {
						continue
					}
					right = res
				case int:
					right = expression.Right.(int)
				}

				uc.service.RequestCalculation(reqId, left, right, expression.Variable, expression.Operation)
			}
			delete(expressionMap, result.Key)
		}

		return false, nil
	}

	err := uc.service.ConsumeResults(reqId, handlerFunc)
	if err != nil {
		return nil, err
	}

	return prints, nil
}
