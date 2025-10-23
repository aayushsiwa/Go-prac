package main

import (
	"context"
	"io"
	"log"

	pb "chatRoom/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	// Open a bidirectional stream
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to open chat stream: %v", err)
	}

	log.Println("👀 Admin spectator connected. Listening for chat messages...\n")

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream ended.")
			return
		}
		if err != nil {
			log.Printf("⚠️ Error receiving message: %v\n", err)
			continue
		}
		log.Printf("[MESSAGE] %s\n", msg.Body)
	}
}
