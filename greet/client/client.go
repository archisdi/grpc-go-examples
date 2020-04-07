package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	doUnary(c, 5*time.Second)
	doUnary(c, 1*time.Second)

	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient, timeout time.Duration) {
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Archie",
			LastName:  "Isdiningrat",
		},
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if response, err := c.Greet(ctx, request); err != nil {

		if resErr, ok := status.FromError(err); ok && resErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Request timeout")
		} else {
			fmt.Println(err.Error())
		}

	} else {
		fmt.Println(response.GetResult())
	}
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

func doBiStreaming(c greetpb.GreetServiceClient) {
	names := []string{
		"angelina",
		"archie",
		"myrtyl",
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error reading stream")
	}

	biChan := make(chan struct{})

	// invoke server
	go func() {
		for _, name := range names {
			stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: name,
					LastName:  "-",
				},
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive
	go func() {
		for {
			if res, err := stream.Recv(); err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error reading server")
				break
			} else {
				fmt.Println(res.GetResult())
			}
		}
		close(biChan)
	}()

	<-biChan
}
