package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "chatRoom/proto"

	"google.golang.org/grpc"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[pb.ChatService_ChatServer]bool
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[pb.ChatService_ChatServer]bool),
	}
}

func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Println("Error receiving message:", err)
			return err
		}

		log.Printf("Received: %s", msg.Body)

		// broadcast message to all clients
		s.mu.Lock()
		for client := range s.clients {
			if err := client.Send(msg); err != nil {
				log.Println("Error sending to client:", err)
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, NewChatServer())

	log.Println("🚀 Chat server started on :9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
