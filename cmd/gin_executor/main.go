package main

import (
	"Calculator/config"
	"Calculator/internal/executor/handlers/gin_handler"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/infrastructure/factories"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"

	_ "Calculator/api/executorswag"
)

// @title	Calculator API
// @version	1.0
func main() {
	cfg := config.LoadYamlConfig("config/config.yaml")

	client := factories.ProduceArithmClient(cfg.ArithmeticServer.Host, cfg.ArithmeticServer.Port)
	broker := factories.ProduceRabbitMQClient(cfg.RabbitMQBroker.URI, 2, 10*time.Second)
	defer broker.Close()

	commService := services.NewCommService(client, broker)
	validService := services.NewValidationService()
	getService := services.NewGetterService()
	useCase := use_cases.NewUseCase(commService, validService, getService)

	handler := gin_handler.NewGinHandler(useCase)
	r := gin.Default()
	r.Use(gin_handler.ReqIdMiddleware())
	r.POST("/execute", handler.Execute)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = r.Run(cfg.ExecutorGinServer.Address + ":" + cfg.ExecutorGinServer.Port)
}
