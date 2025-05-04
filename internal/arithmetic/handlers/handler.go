package handlers

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/use_cases"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type ArithmeticServer struct {
	arithmeticpb.UnimplementedArithmeticServer
	uc *use_cases.UseCase
}

func Register(gs *grpc.Server, uc *use_cases.UseCase) {
	arithmeticpb.RegisterArithmeticServer(gs, &ArithmeticServer{uc: uc})
}

func (as *ArithmeticServer) Calculate(
	ctx context.Context,
	in *arithmeticpb.CalculationData,
) (*arithmeticpb.Message, error) {
	expression := arithmetic.Expression{
		Variable: in.GetVar(),
		Op:       in.GetOp(),
		Left:     in.GetLeft(),
		Right:    in.GetRight(),
	}

	go as.uc.Execute(expression, in.GetQueueName())

	return &arithmeticpb.Message{
		Text: fmt.Sprintf("Expression received: %v", expression),
	}, nil
}
