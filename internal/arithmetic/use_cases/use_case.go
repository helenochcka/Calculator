package use_cases

import (
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/services"
	"fmt"
	"time"
)

type UseCase struct {
	brokerService services.ResultService
	arithmService services.ArithmService
}

func NewUseCase(bs services.ResultService, as services.ArithmService) UseCase {
	return UseCase{brokerService: bs, arithmService: as}
}

func (uc *UseCase) Execute(expression arithmetic.Expression, queueName string) {
	time.Sleep(arithmetic.WorkSimulationTime)

	opToArithmFuncMap := map[string]func(left, right int64) int64{
		"+": uc.arithmService.Sum,
		"-": uc.arithmService.Sub,
		"*": uc.arithmService.Multi,
	}

	arithmFunc, ok := opToArithmFuncMap[expression.Op]
	if !ok {
		msg := fmt.Sprintf("operation '%s' is not supported", expression.Op)
		uc.brokerService.PublishError(msg, queueName)
		return
	}

	result := arithmetic.Result{Key: expression.Variable}
	result.Value = arithmFunc(expression.Left, expression.Right)

	uc.brokerService.PublishResult(result, queueName)
}
