proto1:
	protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.


proto2:
	protoc calculator/calculatorpb/calculator.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.

