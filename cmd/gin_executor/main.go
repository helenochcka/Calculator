package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/gin_handlers"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure/factories"
	"github.com/gin-gonic/gin"
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

	handler := gin_handlers.NewGinHandler(useCase)
	r := gin.Default()
	r.Use(gin_handlers.ReqIdMiddleware())
	r.POST("/execute", handler.Execute)
	_ = r.Run(cfg.ExecutorGinServer.Port)
}
