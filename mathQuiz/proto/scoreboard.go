package proto

import (
	"log"
)

func (s *Server) StreamScoreboard(_ *Empty, stream QuizService_StreamScoreboardServer) error {
	lastSent := 0

	for {
		<-s.scoreboardUpdated

		s.mu.Lock()
		if len(s.scores) > lastSent {
			newScores := s.scores[lastSent:]
			lastSent = len(s.scores)
			s.mu.Unlock()

			for _, ps := range newScores {
				score := &Score{
					PlayerId: ps.ID,
					Correct:  int32(ps.Correct),
					Total:    int32(ps.Total),
				}
				if err := stream.Send(&Scoreboard{Scores: []*Score{score}}); err != nil {
					log.Println("Error streaming scoreboard:", err)
					return err
				}
			}
		} else {
			s.mu.Unlock()
		}
	}
}
