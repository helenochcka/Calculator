package service

import (
	"Calculator/api/arithmeticpb"
	"context"
)

type Service struct {
	client arithmeticpb.ArithmeticServiceClient
	broker *RabbitMQBroker
}

func NewService(c arithmeticpb.ArithmeticServiceClient, b *RabbitMQBroker) Service {
	return Service{client: c, broker: b}
}

func (s *Service) RequestCalculation(left, right int, variable, op string) {
	req := arithmeticpb.CalculationData{
		Literal: variable,
		Op:      op,
		Left:    int64(left),
		Right:   int64(right),
	}
	_, err := s.client.Calculate(context.Background(), &req)

	if err != nil {
		return
	}
}

func (s *Service) DeclareQueue(queueName string) error {
	err := s.broker.DeclareQueue(queueName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ConsumeResults(queue string, handler MessageHandler) error {
	err := s.broker.Consume(queue, handler)
	if err != nil {
		return err
	}
	return nil
}
