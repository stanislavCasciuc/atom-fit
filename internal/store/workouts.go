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
	Likes            int                `json:"likes"`
	Rating           float32            `json:"rating"`
	ReviewsCount     int                `json:"reviews_count"`
	UserLiked        bool               `json:"user_liked"`
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

func (s *WorkoutStore) GetAll(
	ctx context.Context,
	fq pagination.PaginatedQuery,
	userID int64,
) ([]Workout, int, error) {
	query := `
	SELECT 
	  w.id, 
	  w.user_id, 
	  w.name, 
	  w.description, 
	  w.tutorial_link, 
	  w.created_at, 
	  COUNT(DISTINCT wl.user_id) AS likes, 
	  COUNT(DISTINCT wr.user_id) AS reviews_count, 
	  COALESCE(AVG(wr.rating), 0.0) AS average_rating,
	  CASE 
	    WHEN EXISTS (
	      SELECT 1 FROM workout_likes wl2 
	      WHERE wl2.workout_id = w.id AND wl2.user_id = $5
	    ) THEN true
	    ELSE false 
	  END AS user_liked ,
    COUNT(*) OVER() AS total_count
	FROM workouts w 
	LEFT JOIN workout_exercises we ON we.workout_id = w.id 
	LEFT JOIN exercises e ON we.exercise_id = e.id 
	LEFT JOIN workout_likes wl ON w.id = wl.workout_id 
	LEFT JOIN workout_reviews wr ON w.id = wr.workout_id 
	WHERE 
	  (e.name ILIKE '%' || $1 || '%' 
	  OR e.description ILIKE '%' || $1 || '%' 
	  OR w.name ILIKE '%' || $1 || '%' 
	  OR w.description ILIKE '%' || $1 || '%') 
	  AND (e.muscles @> $2 OR $2 = '{}') 
	GROUP BY w.id 
	ORDER BY likes ` + fq.Sort + ` 
	LIMIT $3 OFFSET $4
`

	workouts := make([]Workout, 0)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fq.Search,
		pq.Array(fq.Tags),
		fq.Limit,
		fq.Offset,
		userID,
	)
	if err != nil {
		return nil, 0, err
	}
	var totalCount int
	for rows.Next() {
		var w Workout

		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Name,
			&w.Description,
			&w.TutorialLink,
			&w.CreatedAt,
			&w.Likes,
			&w.ReviewsCount,
			&w.Rating,
			&w.UserLiked,
			&totalCount,
		)
		if err != nil {
			return nil, 0, err
		}

		workouts = append(workouts, w)
	}
	return workouts, totalCount, nil
}

func (s *WorkoutStore) GetUsersWorkouts(
	ctx context.Context,
	fq pagination.PaginatedQuery,
	userID int64,
) ([]Workout, error) {
	query := `
		SELECT 
	  w.id, 
	  w.user_id, 
	  w.name, 
	  w.description, 
	  w.tutorial_link, 
	  w.created_at, 
	  COUNT(DISTINCT wl.user_id) AS likes, 
	  COUNT(DISTINCT wr.user_id) AS reviews_count, 
	  COALESCE(AVG(wr.rating), 0.0) AS average_rating,
	  CASE 
	    WHEN EXISTS (
	      SELECT 1 FROM workout_likes wl2 
	      WHERE wl2.workout_id = w.id AND wl2.user_id = $1
	    ) THEN true
	    ELSE false 
	  END AS user_liked 
	FROM workouts w 
	LEFT JOIN workout_exercises we ON we.workout_id = w.id 
	LEFT JOIN exercises e ON we.exercise_id = e.id 
	LEFT JOIN workout_likes wl ON w.id = wl.workout_id 
	LEFT JOIN workout_reviews wr ON w.id = wr.workout_id 
	WHERE 
	 w.user_id = $1 
	GROUP BY w.id 
	ORDER BY likes ` + fq.Sort + ` 
	LIMIT $2 OFFSET $3

	`
	workouts := make([]Workout, 0)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
		fq.Limit,
		fq.Offset,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var w Workout

		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Name,
			&w.Description,
			&w.TutorialLink,
			&w.CreatedAt,
			&w.Likes,
			&w.ReviewsCount,
			&w.Rating,
			&w.UserLiked,
		)
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

func (s *WorkoutStore) GetByID(ctx context.Context, id int64) (*Workout, error) {
	query := `
		SELECT id, user_id, name, description, tutorial_link, created_at FROM workouts 
		WHERE id = $1
	`
	w := &Workout{}
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&w.ID, &w.UserID, &w.Name, &w.Description, &w.TutorialLink, &w.CreatedAt)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *WorkoutStore) GetWorkoutExercises(
	ctx context.Context,
	workoutID int64,
) ([]WorkoutExercises, error) {
	query := `
		SELECT id, user_id, name, description, is_duration, e.duration, tutorial_link, created_at, muscles, we.duration FROM exercises e 
		JOIN workout_exercises we ON e.id = we.exercise_id
		WHERE we.workout_id = $1
	`
	workoutExercises := make([]WorkoutExercises, 0)
	rows, err := s.db.QueryContext(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e Exercise
		var duration int
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.Name,
			&e.Description,
			&e.IsDuration,
			&e.Duration,
			&e.TutorialLink,
			&e.CreatedAt,
			pq.Array(&e.Muscles),
			&duration,
		)
		if err != nil {
			return nil, err
		}
		we := WorkoutExercises{
			ExerciseID: e.ID,
			Exercise:   e,
			WorkoutID:  workoutID,
			Duration:   duration,
		}

		workoutExercises = append(workoutExercises, we)
	}

	return workoutExercises, nil
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
