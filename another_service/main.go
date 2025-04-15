package main

import (
	grpcServer "Calculator/another_service/grpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	gRPCServer := grpc.NewServer()
	grpcServer.Register(gRPCServer)
	grpcListener, _ := net.Listen("tcp", ":50051")
	gRPCServer.Serve(grpcListener)
}
