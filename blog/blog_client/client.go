package main

import (
	"context"
	"gRPC_project/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("")
	}

	client := blogpb.NewBlogServiceClient(dial)
	defer dial.Close()
	//createBlog(client)
	readBlog(client)

}

func createBlog(c blogpb.BlogServiceClient) {

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Yemi",
			Title:    "Fist Blog",
			Content:  "This is my first blog",
		},
	}
	blog, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			if stat.Code() == codes.Internal {
				log.Fatalf("Internal error: %v", stat.Message())
			} else {
				log.Fatalf("Unexpected error: %v", stat.Message())
			}
		}
		log.Fatalf("Unexpected error: %v", err)
	}

	log.Printf("Blog created: %v", blog)

}

func readBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.ReadBlogRequest{
		Id: "6277e28a18eb6a73db18a892",
	}

	blog, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	log.Printf("Blog read: %v", blog)

}
