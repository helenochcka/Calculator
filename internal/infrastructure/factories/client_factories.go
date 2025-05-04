package factories

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ProduceArithmClient(address, port string) arithmeticpb.ArithmeticClient {
	conn, err := grpc.NewClient(address+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}

	conn.Connect()

	return arithmeticpb.NewArithmeticClient(conn)
}

func ProduceRabbitMQClient(uri, ct string) *rabbitmq.Client {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ broker: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("failed to open AMQP channel: %v", err)
	}

	return rabbitmq.NewClient(conn, ch, ct)
}
