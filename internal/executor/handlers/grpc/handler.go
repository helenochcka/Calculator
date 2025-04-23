package grpc

import (
	"Calculator/api/executorpb"
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_case"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	Instructions := make([]executor.Instruction, 0, len(in.GetInstructions()))
	for _, genInst := range in.GetInstructions() {
		instruction := executor.Instruction{
			Type:      genInst.Type,
			Operation: genInst.Op,
			Result:    genInst.Var,
			Right:     genInst.Right,
			Left:      genInst.Left,
		}
		Instructions = append(Instructions, instruction)
	}

	items, err := s.uc.Execute(Instructions)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
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
