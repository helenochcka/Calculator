package main

import (
	"Calculator/config"
	ginHandler "Calculator/internal/executor/handlers/gin"
	executorService "Calculator/internal/executor/service"
	"Calculator/internal/executor/stores/concurrent_map"
	executorUseCase "Calculator/internal/executor/use_case"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	instructionStorage := concurrent_map.NewInstructionStorage()
	resultStorage := concurrent_map.NewResultStorage()

	client := executorService.ProduceClient()

	service := executorService.NewService(instructionStorage, resultStorage, client)
	useCase := executorUseCase.NewUseCase(service)

	handler := ginHandler.NewHandlerGin(useCase)
	r := gin.Default()
	r.POST("/calculate", handler.Calculate)
	_ = r.Run(cfg.ExecutorServer.Port)
}
