package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "chatRoom/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating chat stream: %v", err)
	}

	// Listen for incoming messages
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Receive error: %v", err)
				continue
			}
			fmt.Printf("\n👤 %s\n> ", msg.Body)
		}
	}()

	// Send user input messages
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("💬 Anonymous Chat Started")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		stream.Send(&pb.ChatMessage{Body: text})
	}
}
