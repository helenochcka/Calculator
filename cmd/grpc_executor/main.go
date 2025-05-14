package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/grpc_handler"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure/factories"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := factories.ProduceArithmClient(cfg.ArithmeticServer.Host, cfg.ArithmeticServer.Port)
	broker := factories.ProduceRabbitMQClient(cfg.RabbitMQBroker.URI, 2, 10*time.Second)
	defer broker.Close()

	commService := services.NewCommService(client, broker)
	validService := services.NewValidationService()
	getService := services.NewGetterService()
	useCase := use_cases.NewUseCase(commService, validService, getService)

	grpcListener, err := net.Listen("tcp", cfg.ExecutorGRPCServer.Address+":"+cfg.ExecutorGRPCServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_handler.ReqIdInterceptor()))
	grpc_handler.Register(server, useCase)
	server.Serve(grpcListener)
}
