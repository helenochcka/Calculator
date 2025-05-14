package factories

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func ProduceArithmClient(host, port string) arithmeticpb.ArithmeticClient {
	conn, err := grpc.NewClient(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc_handler connection: %v", err)
	}

	conn.Connect()

	return arithmeticpb.NewArithmeticClient(conn)
}

func ProduceRabbitMQClient(uri string, maxRetries int, delay time.Duration) *rabbitmq.Client {
	var conn *amqp.Connection
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(uri)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ not ready yet, retrying in %s... (%d/%d)", delay, i+1, maxRetries)
		time.Sleep(delay)
	}

	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ broker: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("failed to open AMQP channel: %v", err)
	}

	return rabbitmq.NewClient(conn, ch)
}
