package grpc_handlers

import (
	"Calculator/internal/executor/values"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func ReqIdInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		reqId := uuid.NewString()
		ctx = context.WithValue(ctx, values.RequestIdKey, reqId)
		return handler(ctx, req)
	}
}
