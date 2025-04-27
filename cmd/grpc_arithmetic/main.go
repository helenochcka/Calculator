package main

import (
	"Calculator/config"
	grpcServer "Calculator/internal/arithmetic"
	executorService "Calculator/internal/executor/service"
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

	broker, err := executorService.NewRabbitMQBroker("amqp://guest:guest@localhost:5672/", "application/x-protobuf")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	defer broker.Close()

	grpcServer.Register(server, broker)
	if err := server.Serve(grpcListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
