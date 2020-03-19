package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"greet/greetpb"
)

func main() {
	fmt.Println("im client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Isdiningrat",
		},
	}
	response, _ := c.Greet(context.Background(), request)
	fmt.Println(response.GetResult())
}
