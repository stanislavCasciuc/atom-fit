package store

import (
	"context"
	"database/sql"
)

type FinishedWorkout struct {
	Date      string
	UserID    int64
	WorkoutID int64
	Duration  string
}

type FinishedWorkoutsStore struct {
	db *sql.DB
}

func (s FinishedWorkoutsStore) Create(ctx context.Context, fn *FinishedWorkout) error {
	query := `
		INSERT INTO finished_workouts (user_id, workout_id, duration) VALUES ($1,$2,$3)
	`
	_, err := s.db.ExecContext(ctx, query, fn.UserID, fn.WorkoutID, fn.Duration)
	if err != nil {
		return err
	}
	return nil
}

func (s FinishedWorkoutsStore) GetAll(
	ctx context.Context,
	userID int64,
) ([]FinishedWorkout, error) {
	query := `
		SELECT date, user_id, workout_id, duration FROM finished_workouts
		WHERE user_id = $1
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	finisedWorkouts := make([]FinishedWorkout, 0)
	for rows.Next() {
		var fn FinishedWorkout
		rows.Scan(
			&fn.Date,
			&fn.UserID,
			&fn.WorkoutID,
			&fn.Duration,
		)
		finisedWorkouts = append(finisedWorkouts, fn)
	}

	return finisedWorkouts, nil
}
