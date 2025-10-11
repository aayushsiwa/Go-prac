package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"gRPC/chat"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	message := &chat.Message{Body: "Hello from Client!"}

	response, err := c.SayHello(context.Background(), message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %v", err)
	}

	log.Printf("Response from server: %s", response.Body)
}
