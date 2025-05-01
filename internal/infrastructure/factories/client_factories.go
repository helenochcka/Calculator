package factories

import (
	"Calculator/api/arithmeticpb"
	"Calculator/config"
	"Calculator/internal/infrastructure/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ProduceArithmClient(cfg config.Config) arithmeticpb.ArithmeticClient {
	conn, err := grpc.NewClient(
		cfg.ArithmeticServer.Address+cfg.ArithmeticServer.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}

	conn.Connect()

	return arithmeticpb.NewArithmeticClient(conn)
}

func ProduceRabbitMQClient(uri string, ct string) *rabbitmq.Client {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ broker: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return rabbitmq.NewClient(conn, ch, ct)
}
