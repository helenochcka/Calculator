package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/grpc_handlers"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure"
	"Calculator/internal/infrastructure/rabbitmq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := infrastructure.ProduceArithmClient(cfg)
	broker := rabbitmq.NewRabbitMQBroker(cfg.RabbitMQBroker.URI, cfg.RabbitMQBroker.ContentType)
	defer broker.Close()

	commService := services.NewCommService(client, broker)
	useCase := use_cases.NewUseCase(commService)

	grpcListener, err := net.Listen("tcp", cfg.ExecutorServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_handlers.ReqIdInterceptor()))
	grpc_handlers.Register(server, useCase)
	server.Serve(grpcListener)
}
