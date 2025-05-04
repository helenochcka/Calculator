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

type ServerAPI struct {
	executorpb.UnimplementedExecutorServer
	uc *use_cases.UseCase
}

func Register(
	gRPCServer *grpc.Server,
	useCase *use_cases.UseCase) {
	executorpb.RegisterExecutorServer(gRPCServer, &ServerAPI{uc: useCase})
}

func (s *ServerAPI) Execute(
	ctx context.Context,
	in *executorpb.Request,
) (*executorpb.Response, error) {

	gi := dto.GroupedInstructions{
		Expressions: make([]executor.Expression, 0),
		VarsToPrint: make(map[string]bool),
	}

	err := s.groupInstructions(in.GetInstructions(), &gi)
	if err != nil {
		return nil, s.gRPCErrMap(err)
	}

	_, ok := ctx.Value(values.RequestIdKey).(string)
	if !ok {
		return nil, status.Error(codes.Internal, errors.New("request id is missing in context").Error())
	}

	results, err := s.uc.Execute(ctx, &gi)
	if err != nil {
		return nil, s.gRPCErrMap(err)
	}

	return &executorpb.Response{Items: *s.resultsToItems(&results)}, nil
}

func (s *ServerAPI) resultsToItems(results *[]executor.Result) *[]*executorpb.Item {
	items := make([]*executorpb.Item, len(*results))
	for i, result := range *results {
		items[i] = &executorpb.Item{
			Var:   result.Key,
			Value: int64(result.Value),
		}
	}
	return &items
}

func (s *ServerAPI) groupInstructions(instructions []*executorpb.Instruction, gi *dto.GroupedInstructions) error {
	for _, instruction := range instructions {
		if instruction.Type == values.Calculate {
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
			continue

		} else if instruction.Type == values.Print {
			gi.VarsToPrint[instruction.Var] = false
			continue
		}
		return fmt.Errorf("%w: %v", executor.ErrUnknownInstructionType, instruction.Type)
	}
	return nil
}

func (s *ServerAPI) gRPCErrMap(err error) error {
	switch {
	case errors.Is(err, executor.ErrCyclicDependency) ||
		errors.Is(err, executor.ErrUnknownInstructionType) ||
		errors.Is(err, executor.ErrCalcExpression) ||
		errors.Is(err, executor.ErrVarNeverBeCalc):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, executor.ErrVarAlreadyUsed):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, executor.ErrVarToPrintNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
