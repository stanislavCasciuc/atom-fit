package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
)

type Workout struct {
	ID               int64              `json:"id"`
	UserID           int64              `json:"user_id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	TutorialLink     string             `json:"tutorial_link"`
	CreatedAt        string             `json:"created_at"`
	WorkoutExercises []WorkoutExercises `json:"workout_exercises"`
}

type WorkoutExercises struct {
	WorkoutID  int64    `json:"workout_id"`
	ExerciseID int64    `json:"exercise_id"`
	Duration   int      `json:"duration"`
	Exercise   Exercise `json:"exercise"`
}

type WorkoutStore struct {
	db *sql.DB
}

// func (s *WorkoutStore) GetWorkoutByID(ctx context.Context, id int64) (*Workout, err) {
// 	query := `
//     SELECT w.id, w.user_id, w.name, w.description, w.tutorial_link, w.created_at, e.id, e.name, e.is_duration, we.duraion, e.description, e.tutorial_link, e.muscles
//     FROM workouts w
//     LEFT JOIN workout_exercises we ON we.workout_id = w.id
//     JOIN exercise e ON we.exercise_id = e.id
//     WHERE w.id = 1
//     GROUP BY w.id
//   `
// 	w := &Workout{}
//
// 	err := s.db.QueryRowContext(ctx, query, id).
// 		Scan(&w.ID, &w.UserID, &w.Name, &w.Description, &w.TutorialLink, &w.CreatedAt, &w.WorkoutExercises.Exercise)
// }

func (s *WorkoutStore) GetAll(
	ctx context.Context,
	fq pagination.PaginatedQuery,
) ([]Workout, error) {
	query := `
	SELECT w.id, w.user_id, w.name, w.description, w.tutorial_link, w.created_at FROM workouts w 
	LEFT JOIN workout_exercises we ON we.workout_id = w.id 
	JOIN exercises e ON we.exercise_id = e.id 
    WHERE (e.name ILIKE '%' || $1 || '%' OR e.description ILIKE '%' || $1 || '%' OR w.name ILIKE '%' || $1 || '%' OR w.description ILIKE '%' || $1 || '%') AND 
(e.muscles @> $2 OR $2 = '{}')
    ORDER BY w.created_at ` + fq.Sort + `
    LIMIT $3 OFFSET $4
  `
	workouts := make([]Workout, 0)

	rows, err := s.db.QueryContext(ctx, query, fq.Search, pq.Array(fq.Tags), fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var w Workout

		err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.Description, &w.TutorialLink, &w.CreatedAt)
		if err != nil {
			return nil, err
		}

		workouts = append(workouts, w)
	}
	return workouts, nil
}

func (s *WorkoutStore) Create(ctx context.Context, w *Workout) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.createWorkout(ctx, tx, w); err != nil {
			return err
		}

		for _, e := range w.WorkoutExercises {
			e.WorkoutID = w.ID
			if err := s.addExerciseToWorkout(ctx, tx, &e); err != nil {
				if errors.Is(err, ErrConflict) {
					continue
				}
				return err
			}
		}
		return nil
	})
}

func (s *WorkoutStore) createWorkout(ctx context.Context, tx *sql.Tx, w *Workout) error {
	query := `
    INSERT INTO workouts (user_id, name, description, tutorial_link) VALUES ($1, $2, $3, $4)
    RETURNING id, created_at
  `
	err := tx.QueryRowContext(ctx, query, w.UserID, w.Name, w.Description, w.TutorialLink).
		Scan(&w.ID, &w.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkoutStore) addExerciseToWorkout(
	ctx context.Context,
	tx *sql.Tx,
	we *WorkoutExercises,
) error {
	query := `
    INSERT INTO workout_exercises (exercise_id, workout_id, duration) VALUES($1, $2, $3)
  `
	_, err := tx.ExecContext(ctx, query, we.ExerciseID, we.WorkoutID, we.Duration)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}
