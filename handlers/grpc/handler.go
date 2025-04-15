package grpc

import (
	"Calculator/core"
	gen "Calculator/handlers/grpc/gen"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	gen.UnimplementedCalcServer
	uc core.UseCase
}

func Register(
	gRPCServer *grpc.Server,
	useCase core.UseCase) {
	gen.RegisterCalcServer(gRPCServer, &serverAPI{uc: useCase})
}

func (s *serverAPI) Calculate(
	ctx context.Context,
	in *gen.Request,
) (*gen.Response, error) {

	Instructions := make([]core.Instruction, 0, len(in.GetInstructions()))
	for _, genInst := range in.GetInstructions() {
		instruction := core.Instruction{
			Type:      genInst.Type,
			Operation: genInst.Op,
			Result:    genInst.Var,
			Right:     genInst.Right,
			Left:      genInst.Left,
		}
		Instructions = append(Instructions, instruction)
	}

	items, err := s.uc.Execute(&ctx, Instructions)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	genItems := make([]*gen.Item, 0, len(items))
	for _, item := range items {
		genItem := gen.Item{
			Var:   item.Var,
			Value: int64(item.Value),
		}
		genItems = append(genItems, &genItem)
	}

	return &gen.Response{Items: genItems}, nil
}
