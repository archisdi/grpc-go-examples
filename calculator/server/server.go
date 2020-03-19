package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"calculator/calcpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calcpb.AddRequest) (*calcpb.AddResponse, error) {
	fmt.Println("add has been called with %v", req)
	return &calcpb.AddResponse{
		Result: req.GetNumberOne() + req.GetNumberTwo(),
	}, nil
}

func main() {
	fmt.Println("im server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf(err.Error())
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculateServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
