package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

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

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)

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

func doServerStreaming(c greetpb.GreetServiceClient) {
	request := &greetpb.GreetManyTimeRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Isdiningrat",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling greet many times")
	}

	for {
		if msg, err := stream.Recv(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error while receiving greet many times")
		} else {
			fmt.Println(msg.GetResult())
		}
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	names := []string{
		"angelina",
		"archie",
		"myrtyl",
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error reading stream")
	}

	// loop over names
	for i, name := range names {
		fmt.Println("sending message number " + strconv.Itoa(i))
		stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: name,
				LastName:  "-",
			},
		})
		time.Sleep(1000 * time.Millisecond)
	}

	if res, err := stream.CloseAndRecv(); err != nil {
		log.Fatalf("error reading server")
	} else {
		fmt.Println(res.GetResult())
	}

}
