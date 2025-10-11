package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	pb "mathQuiz/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Connecting to server...")
	conn, err := grpc.NewClient("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewQuizServiceClient(conn)
	stream, err := client.Play(context.Background())
	if err != nil {
		log.Fatalf("Error creating quiz stream: %v", err)
	}

	header, err := stream.Header()
	if err != nil {
		log.Fatalf("Error getting header: %v", err)
	}
	nSlice := header.Get("n")
	if len(nSlice) == 0 {
		log.Println("Missing question count in header.")
	}
	nStr := nSlice[0]
	nParsed, _ := strconv.Atoi(nStr)
	n := int32(nParsed)
	fmt.Printf("Quiz has %d questions.\n", n)

	reader := bufio.NewReader(os.Stdin)

	for i := int32(0); i < n; i++ {
		res, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving question: %v", err)
		}
		q := res.GetQuestion()
		if q == nil {
			break
		}

		fmt.Printf("Q%d: %s\n", i+1, q)

		fmt.Print("Your answer: ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		// Send answer to server
		if err := stream.Send(&pb.Answer{Answer: answer}); err != nil {
			log.Fatalf("Error sending answer: %v", err)
			break
		}
	}

	stream.CloseSend()
	resp, err := stream.Recv()
	if err != nil {
		log.Fatalf("Error receiving final score: %v", err)
	}
	score := resp.GetScore()
	fmt.Printf("\nFinal Score: %d/%d\n", score.Correct, score.Total)
}
