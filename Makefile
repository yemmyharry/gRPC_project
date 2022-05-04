proto1:
	protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.


proto2:
	protoc calculator/calculatorpb/calculator.proto --go_out=. --go-grpc_out=.

