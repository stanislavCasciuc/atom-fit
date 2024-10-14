package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type LikesStore struct {
	db *sql.DB
}

func (s *LikesStore) CreateExercise(ctx context.Context, userID, exerciseID int64) error {
	stmt := `
    INSERT INTO exercise_likes (user_id, exercise_id) VALUES($1,$2)
  `

	_, err := s.db.ExecContext(ctx, stmt, userID, exerciseID)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

func (s *LikesStore) DeleteExercise(ctx context.Context, userID, exerciseID int64) error {
	stmt := `
    DELETE FROM exercise_likes WHERE user_id = $1 AND exercise_id = $2
  `

	_, err := s.db.ExecContext(ctx, stmt, userID, exerciseID)
	if err != nil {
		return err
	}

	return nil
}

func (s *LikesStore) CreateWorkout(ctx context.Context, userID, workoutID int64) error {
	stmt := `
    INSERT INTO workout_likes (user_id, workout_id) VALUES($1,$2)
  `

	_, err := s.db.ExecContext(ctx, stmt, userID, workoutID)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

func (s *LikesStore) DeleteWorkout(ctx context.Context, userID, workoutID int64) error {
	stmt := `
    DELETE FROM workout_likes WHERE user_id = $1 AND workout_id = $2
  `

	_, err := s.db.ExecContext(ctx, stmt, userID, workoutID)
	if err != nil {
		return err
	}

	return nil
}
