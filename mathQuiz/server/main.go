package main

import (
	"log"
	"net"

	pb "mathQuiz/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterQuizServiceServer(grpcServer, pb.NewQuizServer())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
