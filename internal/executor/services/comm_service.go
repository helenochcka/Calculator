package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/infrastructure/rabbitmq"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
)

type BrokerClient interface {
	DeclareQueue(name string) error
	Publish(queue string, body []byte) error
	Consume(queue string, rp executor.ResultProcessor) error
	Close() error
}

type ArithmeticClient interface {
	Calculate(ctx context.Context, req *arithmeticpb.CalculationData, opts ...grpc.CallOption) (*arithmeticpb.Message, error)
}

type CommunicationService struct {
	ac ArithmeticClient
	bc BrokerClient
}

func NewCommService(ac ArithmeticClient, bc BrokerClient) *CommunicationService {
	return &CommunicationService{ac: ac, bc: bc}
}

func (cs *CommunicationService) RequestCalculation(cd *dto.CalculationData) {
	req := arithmeticpb.CalculationData{
		Var:       cd.Variable,
		Op:        cd.Operation,
		Left:      int64(cd.Left),
		Right:     int64(cd.Right),
		QueueName: cd.QueueName,
	}

	_, _ = cs.ac.Calculate(context.Background(), &req)
}

func (cs *CommunicationService) DeclareResultsQueue(queueName string) error {
	err := cs.bc.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("%wwith a name: %v", executor.ErrDeclaringQueue, queueName)
	}
	return nil
}

func (cs *CommunicationService) ConsumeResults(queue string, rp executor.ResultProcessor) error {
	err := cs.bc.Consume(queue, rp)
	if err != nil {
		switch {
		case errors.Is(err, rabbitmq.ErrCalculatingResult):
			return fmt.Errorf("%w: %v", executor.ErrCalcExpression, err)
		case errors.Is(err, rabbitmq.ErrUnmarshallingMsg) || errors.Is(err, rabbitmq.ErrConsumingMsgs):
			return fmt.Errorf("%w: %v", executor.ErrConsumingResult, err)
		default:
			return err
		}
	}
	return nil
}
