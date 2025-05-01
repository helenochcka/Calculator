package grpc_handlers

import (
	"Calculator/api/executorpb"
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_cases"
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

	var expressions []executor.Expression
	varsToPrint := make(map[string]bool)
	for _, instruction := range in.GetInstructions() {
		expression, err := s.distributeInstructions(instruction, varsToPrint)
		if err != nil {
			return nil, s.gRPCErrMap(err)
		}
		if expression != nil {
			expressions = append(expressions, *expression)
		}
	}

	_, ok := ctx.Value("request_id").(string)
	if !ok {
		return nil, status.Error(codes.Internal, errors.New("request id is missing in context").Error())
	}

	items, err := s.uc.Execute(ctx, &expressions, varsToPrint)
	if err != nil {
		return nil, s.gRPCErrMap(err)
	}

	genItems := make([]*executorpb.Item, 0, len(items))
	for _, item := range items {
		genItem := executorpb.Item{
			Var:   item.Key,
			Value: int64(item.Value),
		}
		genItems = append(genItems, &genItem)
	}

	return &executorpb.Response{Items: genItems}, nil
}

func (s *ServerAPI) distributeInstructions(
	instruction *executorpb.Instruction,
	varsToPrint map[string]bool) (*executor.Expression, error) {
	if instruction.Type == "calc" {
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
		return &expression, nil

	} else if instruction.Type == "print" {
		varsToPrint[instruction.Var] = false
		return nil, nil
	} else {
		return nil, fmt.Errorf("%w: %v", executor.ErrUnknownInstructionType, instruction.Type)
	}
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
