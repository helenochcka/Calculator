package main

import (
	"Calculator/config"
	grpcServer "Calculator/internal/arithmetic"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	grpcListener, err := net.Listen("tcp", cfg.ArithmeticServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	grpcServer.Register(server)
	if err := server.Serve(grpcListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
