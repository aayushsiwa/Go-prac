package main

import (
	"context"
	"fmt"
	"io"
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

	stream, err := client.StreamScoreboard(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error getting scoreboard: %v", err)
	}

	fmt.Println("🏆 Live Scoreboard Updates:")

	allScores := make([]*pb.Score, 0)

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving scoreboard: %v", err)
		}

		// Each update now has just one score
		newScore := update.Scores[0]
		allScores = append(allScores, newScore)

		fmt.Printf("%d. %s — %d/%d\n",
			len(allScores),
			newScore.GetPlayerId(),
			newScore.GetCorrect(),
			newScore.GetTotal())
	}
}
