proto1:
	protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.


proto2:
	protoc calculator/calculatorpb/calculator.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.

proto3:
		protoc blog/blogpb/blog.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.

greet_server:
	go run greet/greet_server/server.go

calculator_server:
	go run calculator/calculator_server/server.go

greet_client:
	go run greet/greet_client/client.go

calculator_client:
	go run calculator/calculator_client/client.go

blog_server:
	go run blog/blog_server/server.go

blog_client:
	go run blog/blog_client/client.go


format:
	go fmt ./...