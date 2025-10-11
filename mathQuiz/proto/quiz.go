package proto

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"
)

type PlayerScore struct {
	ID      string
	Correct int
	Total   int
}

type Server struct {
	UnimplementedQuizServiceServer
	questions []string
	answers   map[string]string
	scores    []PlayerScore
	mu        sync.Mutex
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

	n := len(s.questions)
	md := metadata.Pairs("n", fmt.Sprintf("%d", n))
	stream.SendHeader(md)

	var playerID string
	correct := 0

	for i, q := range qs {
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

		if playerID == "" {
			playerID = ans.GetPlayerId()
			if playerID == "" {
				playerID = fmt.Sprintf("Player%d", rand.Intn(1000))
			}
			log.Printf("New player joined: %s", playerID)
		}

		if ans.GetAnswer() == s.answers[q] {
			correct++
		}

		log.Printf("[%s] Q%d answered: %s | Score: %d/%d",
			playerID, i+1, ans.GetAnswer(), correct, n)
	}

	score := Score{PlayerId: playerID, Correct: int32(correct), Total: int32(n)}

	s.mu.Lock()
	s.scores = append(s.scores, PlayerScore{ID: playerID, Correct: correct, Total: n})
	s.mu.Unlock()

	if err := stream.Send(&QuizResponse{
		Payload: &QuizResponse_Score{Score: &score},
	}); err != nil {
		return err
	}

	log.Printf("%s finished quiz: %d/%d", playerID, correct, n)
	return nil
}

func (s *Server) GetLeaderboard(ctx context.Context, _ *Empty) (*Leaderboard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var scores []*Score
	for _, ps := range s.scores {
		scores = append(scores, &Score{
			PlayerId: ps.ID,
			Correct:  int32(ps.Correct),
			Total:    int32(ps.Total),
		})
	}

	return &Leaderboard{Scores: scores}, nil
}
