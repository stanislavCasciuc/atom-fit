package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
)

var (
	QueryTimeDuration = time.Second * 5
	ErrNotFound       = errors.New("entity not found")
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		AddUserWeight(context.Context, int64, float32) error
		GetUserAttr(context.Context, int64) (*UserAttributes, error)
		UpdateUserWeight(context.Context, int64, float32) error
		GetUserWeight(context.Context, pagination.PaginatedQuery, int64) ([]UserWeightByDate, error)
	}
	Exercises interface {
		Create(context.Context, *Exercise) error
		GetByID(context.Context, int64) (*Exercise, error)
		GetAll(context.Context, pagination.PaginatedQuery) ([]Exercise, error)
		GetUsersExercises(context.Context, pagination.PaginatedQuery, int64) ([]Exercise, error)
	}
	Workouts interface {
		Create(context.Context, *Workout) error
		GetAll(context.Context, pagination.PaginatedQuery, int64) ([]Workout, int, error)
		GetByID(context.Context, int64) (*Workout, error)
		GetUsersWorkouts(context.Context, pagination.PaginatedQuery, int64) ([]Workout, error)
		GetWorkoutExercises(context.Context, int64) ([]WorkoutExercises, error)
	}
	Likes interface {
		CreateExercise(context.Context, int64, int64) error
		DeleteExercise(context.Context, int64, int64) error
		CreateWorkout(context.Context, int64, int64) error
		DeleteWorkout(context.Context, int64, int64) error
	}
	Reviews interface {
		CreateWorkout(context.Context, *WorkoutReview) error
		Get(context.Context, int64) ([]WorkoutReviewWithMetadata, error)
	}
	FinishedWorkouts interface {
		Create(context.Context, *FinishedWorkout) error
		GetAll(context.Context, int64) ([]FinishedWorkout, error)
	}
}

func New(db *sql.DB) Storage {
	return Storage{
		Users:            &UserStore{db},
		Exercises:        &ExerciseStore{db},
		Likes:            &LikesStore{db},
		Workouts:         &WorkoutStore{db},
		Reviews:          &ReviewsStore{db},
		FinishedWorkouts: &FinishedWorkoutsStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
