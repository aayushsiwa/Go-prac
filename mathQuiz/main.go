package main

import (
	"context"
	"fmt"
	"log"

	pb "mathQuiz/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewQuizServiceClient(conn)

	resp, err := client.GetLeaderboard(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error getting leaderboard: %v", err)
	}

	fmt.Println("🏆 Leaderboard:")
	for i, s := range resp.GetScores() {
		fmt.Printf("%d. %s — %d/%d\n", i+1, s.GetPlayerId(), s.GetCorrect(), s.GetTotal())
	}
}
