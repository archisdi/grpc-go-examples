package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetingResponse, error) {
	fmt.Println("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greeting := "Hello " + firstName + " " + lastName

	response := &greetpb.GreetingResponse{
		Result: greeting,
	}

	return response, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimeRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetMany function was invoked with %v", req)

	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimeResponse{
			Result: "Hola " + firstName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("sending message to client number " + strconv.Itoa(i))
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked")
	message := "Hola "

	for {
		if req, err := stream.Recv(); err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: message,
			})
		} else if err != nil {
			return err
		} else {
			firstName := req.GetGreeting().GetFirstName()
			message += firstName + "! "
		}
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked")

	for {
		if req, err := stream.Recv(); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		} else {
			firstName := req.GetGreeting().GetFirstName()
			message := "Holaa " + firstName + " !"

			stream.Send(&greetpb.GreetEveryoneResponse{
				Result: message,
			})
		}
	}
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
