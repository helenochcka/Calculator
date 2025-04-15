package grpc

import (
	grpc2 "Calculator/another_service/grpc/gen"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type calcAPI struct {
	grpc2.UnimplementedCalcServiceServer
}

func Register(gRPCServer *grpc.Server) {
	grpc2.RegisterCalcServiceServer(gRPCServer, &calcAPI{})
}

func (c *calcAPI) Calculate(ctx context.Context, in *grpc2.CalcRequest) (*grpc2.CalcResponse, error) {
	fmt.Print("start ")

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

	fmt.Println("end")
	return &grpc2.CalcResponse{Result: res}, nil
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
