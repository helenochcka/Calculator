package arithmetic

import (
	"Calculator/api/arithmeticpb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type calcAPI struct {
	arithmeticpb.UnimplementedArithmeticServiceServer
}

func Register(gRPCServer *grpc.Server) {
	arithmeticpb.RegisterArithmeticServiceServer(gRPCServer, &calcAPI{})
}

func (c *calcAPI) Calculate(ctx context.Context, in *arithmeticpb.CalculationData) (*arithmeticpb.Result, error) {
	time.Sleep(50 * time.Millisecond)

	var res int64

	op := in.GetOp()
	left := in.GetLeft()
	right := in.GetRight()

	switch op {
	case "+":
		res = c.sum(left, right)
	case "*":
		res = c.multi(left, right)
	case "-":
		res = c.sub(left, right)
	default:
		return nil, status.Error(codes.Internal, "failed to calculate values")
	}

	return &arithmeticpb.Result{Result: res}, nil
}

func (c *calcAPI) sum(left, right int64) int64 {
	return left + right
}

func (c *calcAPI) multi(left, right int64) int64 {
	return left * right
}

func (c *calcAPI) sub(left, right int64) int64 {
	return left - right
}
