package main

import (
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

	c := greetpb.GreetServiceClient(cc)
	fmt.Println("Created client: %f", c)
}
