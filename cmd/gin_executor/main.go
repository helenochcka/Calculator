package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/gin_handlers"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure"
	"Calculator/internal/infrastructure/rabbitmq"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := infrastructure.ProduceArithmClient(cfg)
	broker := rabbitmq.NewRabbitMQBroker(cfg.RabbitMQBroker.URI, cfg.RabbitMQBroker.ContentType)
	defer broker.Close()

	commService := services.NewCommService(client, broker)
	useCase := use_cases.NewUseCase(commService)

	handler := gin_handlers.NewHandlerGin(useCase)
	r := gin.Default()
	r.Use(gin_handlers.ReqIdMiddleware())
	r.POST("/execute", handler.Execute)
	_ = r.Run(cfg.ExecutorServer.Port)
}
