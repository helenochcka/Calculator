package service

import "Calculator/internal/executor"

type MessageHandler func(result executor.Result) (bool, error)

type Broker interface {
	DeclareQueue(name string) error
	Publish(queue string, body []byte) error
	Consume(queue string, handler MessageHandler) error
	Close() error
}
