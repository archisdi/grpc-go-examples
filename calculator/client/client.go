package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"calculator/calcpb"
)

func main() {
	fmt.Println("im client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer cc.Close()

	c := calcpb.NewCalculateServiceClient(cc)
	add(c)
}

func add(c calcpb.CalculateServiceClient) {
	request := &calcpb.AddRequest{
		NumberOne: 10,
		NumberTwo: 69,
	}
	response, _ := c.Add(context.Background(), request)
	fmt.Println(response.GetResult())
}
