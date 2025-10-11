package proto

import (
	"context"
)

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
