package main

import (
	"Calculator/core"
	//handlerGin "Calculator/handlers/gin"
	grpcServer "Calculator/handlers/grpc"
	"Calculator/stores/concurrent_map"

	//"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
)

func main() {
	//cfg := config.LoadYamlConfig("config/config.yaml")

	instructionStorage := concurrent_map.NewInstructionStorage()
	resultStorage := concurrent_map.NewResultStorage()

	client := core.ProduceClient()

	service := core.NewService(instructionStorage, resultStorage, client)
	useCase := core.NewUseCase(service)

	gRPCServer := grpc.NewServer()
	grpcServer.Register(gRPCServer, useCase)
	grpcListener, _ := net.Listen("tcp", ":8080")
	gRPCServer.Serve(grpcListener)

	//ginHandler := handlerGin.NewHandlerGin(useCase)
	//r := gin.Default()
	//r.POST("/calculate", ginHandler.Calculate)
	//_ = r.Run(":8080")
}
