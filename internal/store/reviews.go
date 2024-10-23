package store

import (
	"context"
	"database/sql"
)

type WorkoutReview struct {
	UserID    int64  `json:"user_id"`
	WorkoutID int64  `json:"workout_id"`
	Rating    int    `json:"rating"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type WorkoutReviewWithMetadata struct {
	WorkoutReview
	Username string
}

type ReviewsStore struct {
	db *sql.DB
}

func (s *ReviewsStore) CreateWorkout(ctx context.Context, wr *WorkoutReview) error {
	query := `
    INSERT INTO workout_reviews(user_id, workout_id, rating, title, content) VALUES($1, $2, $3, $4, $5)
    RETURNING created_at
  `
	err := s.db.QueryRowContext(ctx, query, wr.UserID, wr.WorkoutID, wr.Rating, wr.Title, wr.Content).
		Scan(&wr.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewsStore) Get(
	ctx context.Context,
	workoutID int64,
) ([]WorkoutReviewWithMetadata, error) {
	query := `
    SELECT user_id, workout_id, rating, title, content, wr.created_at, u.username FROM workout_reviews wr
    LEFT JOIN  users u ON u.ID = wr.user_id
    WHERE workout_id = $1
  `
	rows, err := s.db.QueryContext(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}

	var workoutReviews []WorkoutReviewWithMetadata
	for rows.Next() {
		var wr WorkoutReviewWithMetadata

		err := rows.Scan(
			&wr.UserID,
			&wr.WorkoutID,
			&wr.Rating,
			&wr.Title,
			&wr.Content,
			&wr.CreatedAt,
			&wr.Username,
		)
		if err != nil {
			return nil, err
		}
		workoutReviews = append(workoutReviews, wr)
	}
	return workoutReviews, nil
}
