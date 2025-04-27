package use_case

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/service"
	"errors"
)

type UseCase struct {
	service service.Service
}

func NewUseCase(s service.Service) UseCase {
	return UseCase{service: s}
}

func (uc *UseCase) Execute(expressions []executor.Expression, varsToPrint map[string]bool) ([]executor.Item, error) {
	expressionMap := make(map[string][]executor.Expression)
	resultMap := make(map[string]int)

	var items []executor.Item

	for _, expression := range expressions {
		left, leftIsInt := expression.Left.(int)
		right, rightIsInt := expression.Right.(int)
		if leftIsInt && rightIsInt {
			uc.service.RequestCalculation(left, right, expression.Variable, expression.Operation)
			continue
		}

		if !leftIsInt {
			if depExpressions := expressionMap[expression.Variable]; depExpressions != nil && uc.cyclicCheck(depExpressions, expression.Left.(string)) {
				return nil, errors.New("cyclic")
			}
			expressionMap[expression.Left.(string)] = append(expressionMap[expression.Left.(string)], expression)
		}
		if !rightIsInt {
			if depExpressions := expressionMap[expression.Variable]; depExpressions != nil && uc.cyclicCheck(depExpressions, expression.Right.(string)) {
				return nil, errors.New("cyclic")
			}
			expressionMap[expression.Right.(string)] = append(expressionMap[expression.Right.(string)], expression)
		}
	}

	queueName := "result"

	if err := uc.service.DeclareQueue(queueName); err != nil {
		return nil, err
	}

	handlerFunc := func(result executor.Result) (stop bool, err error) {
		resultMap[result.Lit] = result.Result
		if expressionMap[result.Lit] != nil {
			for _, expression := range expressionMap[result.Lit] {
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

				uc.service.RequestCalculation(left, right, expression.Variable, expression.Operation)
			}
			delete(expressionMap, result.Lit)
		}
		if varsToPrint[result.Lit] {
			items = append(items, executor.Item{Var: result.Lit, Value: resultMap[result.Lit]})
			delete(varsToPrint, result.Lit)
		}

		if len(varsToPrint) == 0 {
			return true, nil
		}

		return false, nil
	}

	err := uc.service.ConsumeResults(queueName, handlerFunc)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *UseCase) cyclicCheck(deps []executor.Expression, lit string) bool {
	for _, dep := range deps {
		if dep.Variable == lit {
			return true
		}
	}
	return false
}
