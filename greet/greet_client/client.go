package main

import (
	"context"
	"fmt"
	"gRPC_project/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	tls := true

	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Failed to load ca trust certificate %v", err)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	//fmt.Printf("Created client: %f", c)

	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiDiStreaming(c)

	//doUnaryWithDeadline(c, 5*time.Second)
	//doUnaryWithDeadline(c, 1*time.Second)
}

func doUnary(cc greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Yemi",
			LastName:  "Harry",
		},
	}

	res, err := cc.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(cc greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Yemi",
			LastName:  "Harry",
		},
	}

	resStream, err := cc.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %v\n", msg.GetResult())
	}

}

func doClientStreaming(cc greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Yemi",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harry",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Joseph",
			},
		},
	}

	stream, err := cc.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res.GetResult())

}

func doBiDiStreaming(cc greetpb.GreetServiceClient) {

	stream, err := cc.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf(" error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Yemi",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Harry",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Joseph",
			},
		},
	}

	waitChan := make(chan struct{})

	go func() {
		// send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Request sent: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		// receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf(" error while receiving: %v", err)
			}
			fmt.Printf(" Response received: %v\n", res.GetResult())
		}
		close(waitChan)
	}()

	<-waitChan

}

func doUnaryWithDeadline(cc greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{FirstName: "Yemi"},
	}
	deadline, err := cc.GreetWithDeadline(ctx, req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			if stat.Code() == codes.DeadlineExceeded {
				fmt.Printf("Timeout was hit! Deadline was exceeded\n")
			} else {
				fmt.Printf("Unexpected error: %v\n", stat)
			}
		} else {
			log.Fatalf(" error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	fmt.Printf("Response: %v\n", deadline.GetResult())
}
