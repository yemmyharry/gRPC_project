package main

import (
	"context"
	"fmt"
	"gRPC_project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("sum function invoked:  %s", req)
	first_num := req.GetFirstNum()
	second_num := req.GetSecondNum()
	result := first_num + second_num

	return &calculatorpb.SumResponse{
		SumResult: result,
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error listening:  %v", err)
	}
}
