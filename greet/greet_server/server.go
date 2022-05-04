package main

import (
	"context"
	"fmt"
	"gRPC_project/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName

	resp := greetpb.GreetResponse{
		Result: result,
	}

	return &resp, nil

}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
