package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetingResponse, error) {
	fmt.Println("greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greeting := "Hello " + firstName + " " + lastName

	response := &greetpb.GreetingResponse{
		Result: greeting,
	}

	return response, nil
}

func main() {
	fmt.Println("im server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
