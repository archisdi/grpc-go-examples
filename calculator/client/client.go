package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

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
}

func add() {
	request := &calcpb.AddRequest{
		NumberOne: 10,
		NumberTwo: 69,
	}
	response, _ := (*Conn).Add(context.Background(), request)
	fmt.Println(response.GetResult())
}
