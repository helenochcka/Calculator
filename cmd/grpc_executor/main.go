package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/grpc_handlers"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure/factories"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := factories.ProduceArithmClient(cfg.ArithmeticServer.Address, cfg.ArithmeticServer.Port)
	broker := factories.ProduceRabbitMQClient(cfg.RabbitMQBroker.URI, cfg.RabbitMQBroker.ContentType)
	defer broker.Close()

	commService := services.NewCommService(client, broker)
	validService := services.NewValidationService()
	getService := services.NewGetterService()
	useCase := use_cases.NewUseCase(commService, validService, getService)

	grpcListener, err := net.Listen("tcp", cfg.ExecutorGRPCServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_handlers.ReqIdInterceptor()))
	grpc_handlers.Register(server, useCase)
	server.Serve(grpcListener)
}
