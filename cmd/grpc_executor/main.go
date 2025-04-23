package main

import (
	"Calculator/config"
	grpcServer "Calculator/internal/executor/handlers/grpc"
	executorService "Calculator/internal/executor/service"
	"Calculator/internal/executor/stores/concurrent_map"
	executorUseCase "Calculator/internal/executor/use_case"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	instructionStorage := concurrent_map.NewInstructionStorage()
	resultStorage := concurrent_map.NewResultStorage()

	client := executorService.ProduceClient()

	service := executorService.NewService(instructionStorage, resultStorage, client)
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
