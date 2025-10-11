package proto

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc/metadata"
)

type Server struct {
	UnimplementedQuizServiceServer
	questions []string
	answers   map[string]string
}

func NewQuizServer() *Server {
	q := []string{
		"5 + 5 ?",
		"12 - 7 ?",
		"5 * 6 ?",
	}
	a := map[string]string{
		"5 + 5 ?":  "10",
		"12 - 7 ?": "5",
		"5 * 6 ?":  "30",
	}

	return &Server{
		questions: q,
		answers:   a,
	}
}

func (s *Server) Play(stream QuizService_PlayServer) error {
	rand.Seed(time.Now().UnixNano())
	qs := append([]string{}, s.questions...) // copy slice
	rand.Shuffle(len(qs), func(i, j int) { qs[i], qs[j] = qs[j], qs[i] })

	log.Println("New Player Joined the Quiz")

	n := len(s.questions)
	md := metadata.Pairs("n", fmt.Sprintf("%d", n))
	stream.SendHeader(md)

	correct := 0

	for _, q := range qs {
		if err := stream.Send(&QuizResponse{
			Payload: &QuizResponse_Question{Question: &Question{Question: q}},
		}); err != nil {
			return err
		}

		ans, err := stream.Recv()
		if err == io.EOF {
			log.Println("Player disconnected.")
			return nil
		}
		if err != nil {
			log.Println("Error receiving answer:", err)
			return err
		}

		if ans.GetAnswer() == s.answers[q] {
			correct++
		}

		fmt.Printf("Answered: %s | Correct: %d/%d\n", ans.GetAnswer(), correct, n)
	}

	log.Println("All questions answered — sending score to client.")
	if err := stream.Send(&QuizResponse{
		Payload: &QuizResponse_Score{Score: &Score{
			Correct: int32(correct),
			Total:   int32(n),
		}},
	}); err != nil {
		return err
	}
	return nil
}
