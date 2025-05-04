package grpc_handlers

import (
	"Calculator/api/executorpb"
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/executor/values"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type ExecutorServer struct {
	executorpb.UnimplementedExecutorServer
	uc *use_cases.UseCase
}

func Register(gs *grpc.Server, uc *use_cases.UseCase) {
	executorpb.RegisterExecutorServer(gs, &ExecutorServer{uc: uc})
}

func (es *ExecutorServer) Execute(ctx context.Context, in *executorpb.Request) (*executorpb.Response, error) {
	gi := dto.GroupedInstructions{
		Expressions: make([]executor.Expression, 0),
		VarsToPrint: make(map[string]bool),
	}

	err := es.groupInstructions(in.GetInstructions(), &gi)
	if err != nil {
		return nil, es.mapExecutorErrToGRPCErr(err)
	}

	results, err := es.uc.Execute(ctx, &gi)
	if err != nil {
		return nil, es.mapExecutorErrToGRPCErr(err)
	}

	return &executorpb.Response{Items: *es.resultsToItems(&results)}, nil
}

func (es *ExecutorServer) resultsToItems(results *[]executor.Result) *[]*executorpb.Item {
	items := make([]*executorpb.Item, len(*results))
	for i, result := range *results {
		items[i] = &executorpb.Item{
			Var:   result.Key,
			Value: int64(result.Value),
		}
	}
	return &items
}

func (es *ExecutorServer) groupInstructions(instructions []*executorpb.Instruction, gi *dto.GroupedInstructions) error {
	for _, instruction := range instructions {
		if instruction.Type == values.Calculate {
			err := es.validateCalcInst(instruction)
			if err != nil {
				return err
			}
			expression := executor.Expression{
				Type:      instruction.Type,
				Operation: *instruction.Op,
				Variable:  instruction.Var,
			}

			right, err := strconv.Atoi(*instruction.Right)
			if err != nil {
				expression.Right = *instruction.Right
			} else {
				expression.Right = right
			}

			left, err := strconv.Atoi(*instruction.Left)
			if err != nil {
				expression.Left = *instruction.Left
			} else {
				expression.Left = left
			}

			gi.Expressions = append(gi.Expressions, expression)
			continue

		} else if instruction.Type == values.Print {
			gi.VarsToPrint[instruction.Var] = true
			continue
		}
		return fmt.Errorf("%v (%v)", "unknown type of instruction", instruction.Type)
	}
	return nil
}

func (es *ExecutorServer) validateCalcInst(instruction *executorpb.Instruction) error {
	if instruction.Op == nil {
		return errors.New("field 'op' is missing in calculate instruction")
	}
	if instruction.Left == nil {
		return errors.New("field 'left' is missing in calculate instruction")
	}
	if instruction.Right == nil {
		return errors.New("field 'right' is missing in calculate instruction")
	}
	return nil
}

func (es *ExecutorServer) mapExecutorErrToGRPCErr(err error) error {
	switch {
	case errors.Is(err, executor.ErrCyclicDependency) ||
		errors.Is(err, executor.ErrCalcExpression) ||
		errors.Is(err, executor.ErrVarWillNeverBeCalc):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, executor.ErrVarAlreadyUsed):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
