package main

import (
	"Calculator/config"
	grpcServer "Calculator/internal/executor/handlers/grpc"
	executorService "Calculator/internal/executor/service"
	executorUseCase "Calculator/internal/executor/use_case"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := executorService.ProduceClient()
	broker, err := executorService.NewRabbitMQBroker("amqp://guest:guest@localhost:5672/", "application/x-protobuf")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	defer broker.Close()

	service := executorService.NewService(client, broker)
	useCase := executorUseCase.NewUseCase(service)

	grpcListener, err := net.Listen("tcp", cfg.ExecutorServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	grpcServer.Register(server, useCase)
	if err := server.Serve(grpcListener); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
