package arithmetic

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type calcAPI struct {
	arithmeticpb.UnimplementedArithmeticServiceServer
	broker *service.RabbitMQBroker
}

func Register(gRPCServer *grpc.Server, b *service.RabbitMQBroker) {
	arithmeticpb.RegisterArithmeticServiceServer(gRPCServer, &calcAPI{broker: b})
}

func (c *calcAPI) Calculate(ctx context.Context, in *arithmeticpb.CalculationData) (*arithmeticpb.Message, error) {
	go func() {
		time.Sleep(50 * time.Millisecond)

		var res int64

		op := in.GetOp()
		left := in.GetLeft()
		right := in.GetRight()

		switch op {
		case "+":
			res = c.sum(left, right)
		case "*":
			res = c.multi(left, right)
		case "-":
			res = c.sub(left, right)
		default:
			return
		}

		c.sendResult(in.Literal, res)
	}()

	return &arithmeticpb.Message{Msg: "Variable will send in queue"}, nil
}

func (c *calcAPI) sum(left, right int64) int64 {
	return left + right
}

func (c *calcAPI) multi(left, right int64) int64 {
	return left * right
}

func (c *calcAPI) sub(left, right int64) int64 {
	return left - right
}

func (c *calcAPI) sendResult(lit string, result int64) {

	queueName := "result"

	if err := c.broker.DeclareQueue(queueName); err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
	}

	msg := &arithmeticpb.Result{Lit: lit, Result: result}

	body, err := proto.Marshal(msg)
	if err != nil {
		log.Panicf("Failed to marshal result: %s", err)
	}

	err = c.broker.Publish(queueName, body)
	if err != nil {
		log.Panicf("Failed to publish a message: %s", err)
	}
	log.Printf("Successfully published a message %s = %d\n", lit, result)
}
