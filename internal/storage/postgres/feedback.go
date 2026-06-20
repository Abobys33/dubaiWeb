package postgres

import (
	"context"
	"fmt"
	"time"
)

type FeedbackRequest struct {
	Name  string
	Phone string
}

func (s *Storage) CreateFeedbackRequest(ctx context.Context, req FeedbackRequest) (int64, error) {
	const query = `
		INSERT INTO feedback_requests (name, phone)
		VALUES ($1, $2)
		RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var id int64
	if err := s.DB.QueryRowContext(ctx, query, req.Name, req.Phone).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create feedback request: %w", err)
	}

	return id, nil
}
