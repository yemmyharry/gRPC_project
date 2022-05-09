package main

import (
	"context"
	"fmt"
	"gRPC_project/calculator/calculatorpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"net/http"
)

type server struct{}

func (s server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Square root function invoked %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Number %v is less than 0", number),
		)
	}
	result := math.Sqrt(float64(number))
	res := &calculatorpb.SquareRootResponse{SquareRoot: float32(result)}
	return res, nil
}

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

	go func() {
		mux := runtime.NewServeMux()
		calculatorpb.RegisterCalculatorServiceHandlerServer(context.Background(), mux, server{})
		log.Fatalln(http.ListenAndServe(":8080", mux))
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error listening:  %v", err)
	}
}
