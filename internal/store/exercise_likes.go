package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type ExercisesLikes struct {
	UserID     int64
	ExerciseID int64
}

type ExercisesLikesStore struct {
	db *sql.DB
}

func (s *ExercisesLikesStore) Create(ctx context.Context, userID, exerciseID int64) error {
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

func (s *ExercisesLikesStore) Delete(ctx context.Context, userID, exerciseID int64) error {
	stmt := `
    DELETE FROM exercise_likes WHERE user_id = $1 AND exercise_id = $2
  `

	_, err := s.db.ExecContext(ctx, stmt, userID, exerciseID)
	if err != nil {
		return err
	}

	return nil
}
