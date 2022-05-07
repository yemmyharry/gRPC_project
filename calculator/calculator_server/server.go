package main

import (
	"context"
	"fmt"
	"gRPC_project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (s server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf(" FindMaximum function was invoked with streaming request \n")
	max := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf(" error while reading client stream: %v", err)
		}
		num := req.GetNumber()

		if num > max {
			max = num
			err := stream.Send(&calculatorpb.FindMaximumResponse{Result: max})
			if err != nil {
				log.Fatalf("Error sending data to client %v\n", err)
			}
		}

	}

}

func (s server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with a streaming request")
	var average int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Average: float32(average)})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		average = req.GetFirstNum() + req.GetSecondNum()/2
	}
}

func (s server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	number := req.GetNumber()
	divisor := int32(2)

	for number > 1 {
		if number%divisor == 0 {
			number = number / divisor
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{PrimeNumbers: divisor})

		} else {
			divisor++
		}
	}
	return nil
}

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
