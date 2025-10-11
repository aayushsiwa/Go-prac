package proto

import (
	"log"
	"time"
)

func (s *Server) StreamLeaderboard(_ *Empty, stream QuizService_StreamLeaderboardServer) error {

	for {
		s.mu.Lock()
		var leaderboard Leaderboard
		for _, ps := range s.scores {
			leaderboard.Scores = append(leaderboard.Scores, &Score{
				PlayerId: ps.ID,
				Correct:  int32(ps.Correct),
				Total:    int32(ps.Total),
			})
		}
		s.mu.Unlock()

		if err := stream.Send(&leaderboard); err != nil {
			log.Println("Error streaming leaderboard:", err)
			return err
		}
		time.Sleep(2 * time.Second)
	}
}
