package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	"calculator/calcpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calcpb.AddRequest) (*calcpb.AddResponse, error) {
	fmt.Println("add has been called with %v", req)
	return &calcpb.AddResponse{
		Result: req.GetNumberOne() + req.GetNumberTwo(),
	}, nil
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Println("square root has been called with %v", req)

	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"received a negative number",
		)
	}

	return &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
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

	// register reflection on gRpc server
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
