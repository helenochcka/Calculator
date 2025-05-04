package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/executor"
	"google.golang.org/protobuf/proto"
	"log"
)

type BrokerClient interface {
	DeclareQueue(name string) error
	Publish(queue string, body []byte) error
	Consume(queue string, handler executor.ResultProcessor) error
	Close() error
}

type ResultService struct {
	bc BrokerClient
}

func NewResultService(bc BrokerClient) *ResultService {
	return &ResultService{bc: bc}
}

func (rs *ResultService) PublishResult(result arithmetic.Result, queueName string) {
	msg := arithmeticpb.Result{Key: &result.Key, Value: &result.Value}

	body, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatalf("failed to marshal result: %s", err)
	}

	err = rs.bc.Publish(queueName, body)
	if err != nil {
		log.Fatalf("failed to publish message: %s", err)
	}
}

func (rs *ResultService) PublishError(errMsg string, queueName string) {
	msg := arithmeticpb.Result{ErrMsg: &errMsg}

	body, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatalf("failed to marshal result: %s", err)
	}

	err = rs.bc.Publish(queueName, body)
	if err != nil {
		log.Fatalf("failed to publish message: %s", err)
	}
}
