package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"calculator/calcpb"
)

var Conn *calcpb.CalculateServiceClient

func main() {
	fmt.Println("im client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer cc.Close()
	connection := calcpb.NewCalculateServiceClient(cc)
	Conn = &connection
	add()
	sqrt()
}

func add() {
	request := &calcpb.AddRequest{
		NumberOne: 10,
		NumberTwo: 69,
	}
	response, _ := (*Conn).Add(context.Background(), request)
	fmt.Println(response.GetResult())
}

func sqrt() {
	requestA := &calcpb.SquareRootRequest{
		Number: -10,
	}

	if response, err := (*Conn).SquareRoot(context.Background(), requestA); err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(response.GetNumberRoot())
	}

}
