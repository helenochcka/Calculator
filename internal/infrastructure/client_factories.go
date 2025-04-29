package infrastructure

import (
	"Calculator/api/arithmeticpb"
	"Calculator/config"
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
