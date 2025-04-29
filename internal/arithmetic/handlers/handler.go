package handlers

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/use_cases"
	"context"
	"google.golang.org/grpc"
)

type ArithmServer struct {
	arithmeticpb.UnimplementedArithmeticServer
	useCase use_cases.UseCase
}

func Register(grpcServer *grpc.Server, uc use_cases.UseCase) {
	arithmeticpb.RegisterArithmeticServer(grpcServer, &ArithmServer{useCase: uc})
}

func (as *ArithmServer) Calculate(c context.Context, in *arithmeticpb.CalculationData) (*arithmeticpb.Message, error) {
	expression := arithmetic.Expression{
		Variable: in.GetVar(),
		Op:       in.GetOp(),
		Left:     in.GetLeft(),
		Right:    in.GetRight(),
	}

	go as.useCase.Execute(expression, in.GetQueueName())

	return &arithmeticpb.Message{
		Text: "Expression received by arithmetic service",
	}, nil
}
