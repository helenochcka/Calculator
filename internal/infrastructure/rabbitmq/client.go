package rabbitmq

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewClient(conn *amqp.Connection, ch *amqp.Channel) *Client {
	return &Client{conn, ch}
}

func (b *Client) DeclareQueue(name string) error {
	_, err := b.ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *Client) Publish(queue string, body []byte) error {
	return b.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/x-protobuf",
			Body:        body,
		},
	)
}

func (b *Client) Consume(queue string, rp executor.ResultProcessor) error {
	msgs, err := b.ch.Consume(
		queue, "", true, false, false, false, nil,
	)
	if err != nil {
		return ErrConsumingMsgs
	}

	for msg := range msgs {
		var pbResult arithmeticpb.Result
		err = proto.Unmarshal(msg.Body, &pbResult)
		if err != nil {
			return ErrUnmarshallingMsg
		}
		if pbResult.ErrMsg != nil {
			return fmt.Errorf("%w%v", ErrCalculatingResult, *pbResult.ErrMsg)
		}
		result := executor.Result{Key: *pbResult.Key, Value: int(*pbResult.Value)}
		stop, err := rp(result)
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

func (b *Client) Close() error {
	if err := b.ch.Close(); err != nil {
		return err
	}
	return b.conn.Close()
}
