package main

import (
	"context"
	"fmt"
	"gRPC_project/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	//fmt.Printf("Created client: %f", c)

	doUnary(c)

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
		log.Fatalf("error while calling Gret RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
