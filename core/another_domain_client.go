package core

import (
	gen "Calculator/another_service/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ProduceClient() gen.CalcServiceClient {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}

	conn.Connect()

	return gen.NewCalcServiceClient(conn)
}
