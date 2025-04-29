package rabbitmq

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"log"
)

type RabbitMQBroker struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	contentType string
}

func NewRabbitMQBroker(uri string, ct string) *RabbitMQBroker {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return &RabbitMQBroker{conn: conn, ch: ch, contentType: ct}
}

func (b *RabbitMQBroker) DeclareQueue(name string) error {
	_, err := b.ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitMQBroker) Publish(queue string, body []byte) error {
	return b.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: b.contentType,
			Body:        body,
		},
	)
}

func (b *RabbitMQBroker) Consume(queue string, handler executor.MessageHandler) error {
	msgs, err := b.ch.Consume(
		queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return FailedConsumeMsgs
	}

	for msg := range msgs {
		var result arithmeticpb.Result
		err := proto.Unmarshal(msg.Body, &result)
		if err != nil {
			return FailedUnmarshalMsg
		}
		if result.ErrMsg != nil {
			return fmt.Errorf("%w%v", UnsuccessfulResult, *result.ErrMsg)
		}
		res := executor.Result{Key: *result.Key, Value: int(*result.Value)}
		stop, err := handler(res)
		if err != nil {
			return err
		}
		if stop {
			b.ch.Cancel("", false)
			break
		}
	}
	return nil
}

func (b *RabbitMQBroker) Close() error {
	if err := b.ch.Close(); err != nil {
		return err
	}
	return b.conn.Close()
}
