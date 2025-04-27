package grpc

import (
	"Calculator/api/executorpb"
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_case"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type serverAPI struct {
	executorpb.UnimplementedExecuteServer
	uc use_case.UseCase
}

func Register(
	gRPCServer *grpc.Server,
	useCase use_case.UseCase) {
	executorpb.RegisterExecuteServer(gRPCServer, &serverAPI{uc: useCase})
}

func (s *serverAPI) Calculate(
	ctx context.Context,
	in *executorpb.Request,
) (*executorpb.Response, error) {

	var expressions []executor.Expression
	varsToPrint := make(map[string]bool)
	for _, instruction := range in.GetInstructions() {
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

			expressions = append(expressions, expression)
		} else if instruction.Type == "print" {
			varsToPrint[instruction.Var] = true
		} else {
			return nil, status.Error(codes.InvalidArgument, "unknown type")
		}
	}

	items, err := s.uc.Execute(expressions, varsToPrint)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed")
	}

	genItems := make([]*executorpb.Item, 0, len(items))
	for _, item := range items {
		genItem := executorpb.Item{
			Var:   item.Var,
			Value: int64(item.Value),
		}
		genItems = append(genItems, &genItem)
	}

	return &executorpb.Response{Items: genItems}, nil
}
