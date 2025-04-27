package main

import (
	"Calculator/config"
	ginHandler "Calculator/internal/executor/handlers/gin"
	executorService "Calculator/internal/executor/service"
	executorUseCase "Calculator/internal/executor/use_case"
	"log"

	"github.com/gin-gonic/gin"
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

	handler := ginHandler.NewHandlerGin(useCase)
	r := gin.Default()
	r.POST("/calculate", handler.Calculate)
	_ = r.Run(cfg.ExecutorServer.Port)
}
