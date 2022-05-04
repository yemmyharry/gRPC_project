package main

import (
	"context"
	"gRPC_project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	cc := calculatorpb.NewCalculatorServiceClient(conn)
	sumUnary(cc)

}

func sumUnary(cc calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNum:  50,
		SecondNum: 50,
	}

	res, err := cc.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling sum RPC %v", err)
	}
	log.Println("Response from Sum: ", res)

}