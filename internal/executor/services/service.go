package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"Calculator/internal/infrastructure/rabbitmq"
	"context"
	"errors"
	"fmt"
	"log"
)

type BrokerClient interface {
	DeclareQueue(name string) error
	Publish(queue string, body []byte) error
	Consume(queue string, handler executor.MessageHandler) error
	Close() error
}

type CommService struct {
	arithmClient arithmeticpb.ArithmeticClient
	brokerClient BrokerClient
}

func NewCommService(ac arithmeticpb.ArithmeticClient, bc BrokerClient) *CommService {
	return &CommService{arithmClient: ac, brokerClient: bc}
}

func (cs *CommService) RequestCalculation(queueName string, left, right int, variable, op string) {
	req := arithmeticpb.CalculationData{
		Var:       variable,
		Op:        op,
		Left:      int64(left),
		Right:     int64(right),
		QueueName: queueName,
	}

	msg, _ := cs.arithmClient.Calculate(context.Background(), &req)

	log.Println(msg)
}

func (cs *CommService) DeclareQueue(queueName string) error {
	err := cs.brokerClient.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("%wwith a name: %v", executor.ErrDeclaringQueue, queueName)
	}
	return nil
}

func (cs *CommService) ConsumeResults(queue string, handler executor.MessageHandler) error {
	err := cs.brokerClient.Consume(queue, handler)
	if err != nil {
		switch {
		case errors.Is(err, rabbitmq.UnsuccessfulResult):
			return fmt.Errorf("%w: %v", executor.ErrCalcExpression, err)
		case errors.Is(err, rabbitmq.FailedUnmarshalMsg) ||
			errors.Is(err, rabbitmq.FailedConsumeMsgs):
			return fmt.Errorf("%w: %v", executor.ErrConsumingResult, err)
		default:
			return err
		}
	}
	return nil
}
