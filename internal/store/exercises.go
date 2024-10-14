package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"

	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
)

type Exercise struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	IsDuration   bool     `json:"is_duration"`
	Duration     int      `json:"duration"`
	TutorialLink string   `json:"tutorial_link"`
	CreatedAt    string   `json:"created_at"`
	Muscles      []string `json:"muscles"`
	Likes        int      `json:"like"`
}

type ExerciseStore struct {
	db *sql.DB
}

func (s *ExerciseStore) GetAll(
	ctx context.Context,
	fq pagination.PaginatedQuery,
) ([]Exercise, error) {
	query := `
		SELECT id, exercises.user_id, name, description, is_duration, duration, tutorial_link, created_at, muscles, COUNT(DISTINCT el.user_id) as likes
    FROM exercises 
		LEFT JOIN exercise_likes el ON exercises.id = el.exercise_id
    WHERE (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%') AND 
(muscles @> $2 OR $2 = '{}')
		GROUP BY id
    ORDER BY likes ` + fq.Sort + `
    LIMIT $3 OFFSET $4
	`

	rows, err := s.db.QueryContext(ctx, query, fq.Search, pq.Array(fq.Tags), fq.Limit, fq.Offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var exercises []Exercise

	for rows.Next() {
		var e Exercise
		log.Println(rows)
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
			&e.Likes,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

func (s *ExerciseStore) Create(ctx context.Context, e *Exercise) error {
	query := ` 
		INSERT INTO exercises (user_id, name, description, is_duration, duration, tutorial_link, muscles) VALUES ($1, $2, $3,$4,$5,$6,$7) 
		RETURNING id, created_at
	`
	err := s.db.QueryRowContext(ctx, query, e.UserID, e.Name, e.Description, e.IsDuration, e.Duration, e.TutorialLink, pq.Array(e.Muscles)).
		Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExerciseStore) GetByID(ctx context.Context, id int64) (*Exercise, error) {
	query := `
		SELECT user_id, name, description, is_duration, duration, tutorial_link, muscles, created_at FROM exercises WHERE id = $1
	`
	e := &Exercise{
		ID: id,
	}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&e.UserID, &e.Name, &e.Description, &e.IsDuration, &e.Duration, &e.TutorialLink, pq.Array(&e.Muscles), &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e, nil
}
