package store

import (
	"context"
	"database/sql"
)

type WorkoutReview struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	WorkoutID int64  `json:"workout_id"`
	Rating    int    `json:"rating"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ReviewsStore struct {
	db *sql.DB
}

func (s *ReviewsStore) CreateWorkout(ctx context.Context, wr *WorkoutReview) error {
	query := `
    INSERT INTO workout_reviews(user_id, workout_id, rating, title, content) VALUES($1, $2, $3, $4, $5)
    RETURNING id, created_at
  `
	err := s.db.QueryRowContext(ctx, query, wr.UserID, wr.WorkoutID, wr.Rating, wr.Title, wr.Content).
		Scan(&wr.ID, &wr.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
