package grpc_handlers

import (
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
		ctx = context.WithValue(ctx, "request_id", reqId)
		return handler(ctx, req)
	}
}
