package service

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"log"
)

type RabbitMQBroker struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	contentType string
}

func NewRabbitMQBroker(url string, ct string) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQBroker{conn: conn, ch: ch, contentType: ct}, nil
}

func (b *RabbitMQBroker) DeclareQueue(name string) error {
	_, err := b.ch.QueueDeclare(name, true, false, false, false, nil)
	return err
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

func (b *RabbitMQBroker) Consume(queue string, handler MessageHandler) error {
	msgs, err := b.ch.Consume(
		queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		var result arithmeticpb.Result
		err := proto.Unmarshal(d.Body, &result)
		if err != nil {
			return err
		}
		res := executor.Result{Lit: result.Lit, Result: int(result.Result)}
		stop, err := handler(res)
		if err != nil {
			log.Printf("Ошибка обработки сообщения: %v", err)
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
