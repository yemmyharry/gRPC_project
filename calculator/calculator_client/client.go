package main

import (
	"context"
	"fmt"
	"gRPC_project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	cc := calculatorpb.NewCalculatorServiceClient(conn)
	//sumUnary(cc)

	//primeFactor(cc)

	//computeAverage(cc)

	findMaximum(cc)
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

func primeFactor(cc calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}

	resStream, err := cc.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("Response from RPC: %v\n ", res)
	}

}

func computeAverage(cc calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			FirstNum:  2,
			SecondNum: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			FirstNum:  4,
			SecondNum: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			FirstNum:  6,
			SecondNum: 4,
		},
	}

	stream, err := cc.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage RPC %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage RPC %v", err)
	}
	fmt.Printf("Response from ComputeAverage: %v\n", res.GetAverage())

}

func findMaximum(cc calculatorpb.CalculatorServiceClient) {
	stream, err := cc.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf(" Error while calling FindMaximum RPC %v", err)
	}

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 7,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 19,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 32,
		},
	}

	waitChan := make(chan struct{})

	go func() {
		//requests := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, req := range requests {
			fmt.Printf("Sending number: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response from FindMaximum RPC %v", err)
			}
			fmt.Printf("New Maximum: %v\n", recv.GetResult())
		}
		close(waitChan)
	}()

	<-waitChan

}
