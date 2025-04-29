package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/infrastructure/rabbitmq"
	"google.golang.org/protobuf/proto"
	"log"
)

type ResultService struct {
	broker *rabbitmq.RabbitMQBroker
}

func NewResultService(b *rabbitmq.RabbitMQBroker) ResultService {
	return ResultService{broker: b}
}

func (rs *ResultService) PublishResult(result arithmetic.Result, queueName string) {
	msg := arithmeticpb.Result{Key: &result.Key, Value: &result.Value}

	body, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatalf("Failed to marshal result: %s", err)
	}

	err = rs.broker.Publish(queueName, body)
	if err != nil {
		log.Fatalf("Failed to publish message: %s", err)
	}
	log.Printf("Message successfully published %s = %d\n", result.Key, result.Value)
}

func (rs *ResultService) PublishError(errMsg string, queueName string) {
	msg := arithmeticpb.Result{ErrMsg: &errMsg}

	body, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatalf("Failed to marshal result: %s", err)
	}

	err = rs.broker.Publish(queueName, body)
	if err != nil {
		log.Fatalf("Failed to publish message: %s", err)
	}
}
