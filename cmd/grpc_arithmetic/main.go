package main

import (
	"Calculator/config"
	"Calculator/internal/arithmetic/handlers"
	"Calculator/internal/arithmetic/services"
	"Calculator/internal/arithmetic/use_cases"
	"Calculator/internal/infrastructure/factories"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	broker := factories.ProduceRabbitMQClient(cfg.RabbitMQBroker.URI, 2, 10*time.Second)
	defer broker.Close()

	resultService := services.NewResultService(broker)
	arithmService := services.NewArithmeticService()
	useCase := use_cases.NewUseCase(resultService, arithmService)

	grpcListener, err := net.Listen("tcp", cfg.ArithmeticServer.Address+":"+cfg.ArithmeticServer.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	handlers.Register(server, useCase)
	server.Serve(grpcListener)
}
