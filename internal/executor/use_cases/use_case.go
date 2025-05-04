package use_cases

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/services"
	"context"
)

type UseCase struct {
	cs *services.CommunicationService
	vs *services.ValidationService
	gs *services.GetterService
}

func NewUseCase(
	cs *services.CommunicationService,
	vs *services.ValidationService,
	gs *services.GetterService,
) *UseCase {
	return &UseCase{cs: cs, vs: vs, gs: gs}
}

func (uc *UseCase) Execute(ctx context.Context, gi *dto.GroupedInstructions) ([]executor.Result, error) {
	reqId, err := uc.gs.GetReqIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if err = uc.cs.DeclareResultsQueue(*reqId); err != nil {
		return nil, err
	}

	dependencyMap := make(map[string][]executor.Expression)
	expressionVars := make(map[string]bool)

	for _, expression := range gi.Expressions {
		if err = uc.vs.CheckIfVarAlreadyUsed(expressionVars, expression.Variable); err != nil {
			return nil, err
		}
		expressionVars[expression.Variable] = true

		if err = uc.vs.ValidateArgType(expression.Left); err != nil {
			return nil, err
		}
		if err = uc.vs.ValidateArgType(expression.Right); err != nil {
			return nil, err
		}

		left, leftIsInt := expression.Left.(int)
		right, rightIsInt := expression.Right.(int)
		if leftIsInt && rightIsInt {

			cd := dto.CalculationData{
				Variable:  expression.Variable,
				Operation: expression.Operation,
				Left:      left,
				Right:     right,
				QueueName: *reqId,
			}
			if err = uc.cs.RequestCalculation(&cd); err != nil {
				return nil, err
			}
			continue
		}

		if !leftIsInt {
			err = uc.vs.CheckCyclicDependency(dependencyMap[expression.Variable], expression.Left.(string))
			if err != nil {
				return nil, err
			}
			dependencyMap[expression.Left.(string)] = append(dependencyMap[expression.Left.(string)], expression)
		}
		if !rightIsInt {
			err = uc.vs.CheckCyclicDependency(dependencyMap[expression.Variable], expression.Right.(string))
			if err != nil {
				return nil, err
			}
			dependencyMap[expression.Right.(string)] = append(dependencyMap[expression.Right.(string)], expression)
		}
	}

	if err = uc.vs.CheckIfArgNeverCalculated(dependencyMap, expressionVars); err != nil {
		return nil, err
	}

	if err = uc.vs.CheckIfPrintVarNeverCalculated(gi.VarsToPrint, expressionVars); err != nil {
		return nil, err
	}

	resultsToPrint := make([]executor.Result, 0, len(gi.VarsToPrint))
	resultMap := make(map[string]int, len(expressionVars))

	resultProcessor := func(result executor.Result) (stop bool, err error) {
		resultMap[result.Key] = result.Value

		if gi.VarsToPrint[result.Key] {
			resultsToPrint = append(resultsToPrint, result)
			delete(gi.VarsToPrint, result.Key)
		}

		if len(gi.VarsToPrint) == 0 {
			return true, nil
		}

		if dependencyMap[result.Key] != nil {
			for _, expression := range dependencyMap[result.Key] {

				left, ok := uc.gs.GetVarValue(expression.Left, resultMap)
				if !ok {
					continue
				}
				right, ok := uc.gs.GetVarValue(expression.Right, resultMap)
				if !ok {
					continue
				}

				cd := dto.CalculationData{
					Variable:  expression.Variable,
					Operation: expression.Operation,
					Left:      *left,
					Right:     *right,
					QueueName: *reqId,
				}
				if err = uc.cs.RequestCalculation(&cd); err != nil {
					return true, err
				}
			}
			delete(dependencyMap, result.Key)
		}

		return false, nil
	}

	err = uc.cs.ConsumeResults(*reqId, resultProcessor)
	if err != nil {
		return nil, err
	}

	return resultsToPrint, nil
}
