package use_cases

import (
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/services"
	"fmt"
	"time"
)

type ResultService interface {
	PublishResult(result arithmetic.Result, queueName string)
	PublishError(errMsg string, queueName string)
}

type UseCase struct {
	rs ResultService
	as *services.ArithmeticService
}

func NewUseCase(rs ResultService, as *services.ArithmeticService) *UseCase {
	return &UseCase{rs: rs, as: as}
}

func (uc *UseCase) Execute(expression arithmetic.Expression, queueName string) {
	time.Sleep(arithmetic.WorkSimulationTime)

	opToArithmFuncMap := map[string]func(left, right int64) int64{
		"+": uc.as.Sum,
		"-": uc.as.Sub,
		"*": uc.as.Multi,
	}

	arithmFunc, ok := opToArithmFuncMap[expression.Op]
	if !ok {
		msg := fmt.Sprintf("operation '%s' is not supported", expression.Op)
		uc.rs.PublishError(msg, queueName)
		return
	}

	result := arithmetic.Result{Key: expression.Variable}
	result.Value = arithmFunc(expression.Left, expression.Right)

	uc.rs.PublishResult(result, queueName)
}
