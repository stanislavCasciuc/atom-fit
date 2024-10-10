package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
)

type UserAttributes struct {
	UserID     int64   `json:"user_id"`
	IsMale     bool    `json:"is_male"`
	Height     int     `json:"height"`
	Goal       string  `json:"goal"        validate:"required,oneof=lose gain maintain"`
	WeightGoal float32 `json:"weight_goal"`
	Weight     float32 `json:"weight"`
	Age        int64   `json:"age"`
}

type UserWeightByDate struct {
	Date   string  `json:"date"`
	Weight float32 `json:"weight"`
}

var ErrConflict = errors.New("conflict violation")

func (s *UserStore) AddUserWeight(
	ctx context.Context,
	userID int64,
	weight float32,
) error {
	query := `
		INSERT INTO user_weight (user_id, weight) VALUES ($1, $2)
	`

	_, err := s.db.ExecContext(ctx, query, userID, weight)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserWeight(
	ctx context.Context,
	userID int64,
	weight float32,
) error {
	query := `
		UPDATE user_weight SET weight = $1 WHERE user_id = $2 AND date = CURRENT_DATE
	`

	_, err := s.db.ExecContext(ctx, query, weight, userID)
	if err != nil {
		return err
	}

	return nil
}

// this return WITH last logged weight
func (s *UserStore) GetUserAttr(ctx context.Context, userID int64) (*UserAttributes, error) {
	query := `
	SELECT ua.user_id, ua.is_male, ua.height, ua.goal, ua.weight_goal, uw.weight, ua.age 
	FROM user_attributes ua 
	JOIN user_weight uw ON ua.user_id = uw.user_id 
	WHERE ua.user_id = $1 AND uw.user_id = $1 
	ORDER BY uw.date DESC
	LIMIT 1
	`
	userAttr := &UserAttributes{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&userAttr.UserID,
		&userAttr.IsMale,
		&userAttr.Height,
		&userAttr.Goal,
		&userAttr.WeightGoal,
		&userAttr.Weight,
		&userAttr.Age)
	if err != nil {
		return nil, err
	}

	return userAttr, nil
}

func (s *UserStore) GetUserWeight(
	ctx context.Context,
	fq pagination.PaginatedQuery,
	userID int64,
) ([]UserWeightByDate, error) {
	query := `
	SELECT date, weight FROM user_weight
	WHERE user_ID = $1
	LIMIT $2 OFFSET $3
  `

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	var userWeight []UserWeightByDate
	for rows.Next() {
		var uw UserWeightByDate
		err := rows.Scan(&uw.Date, &uw.Weight)
		if err != nil {
			return nil, err
		}
		userWeight = append(userWeight, uw)

	}
	return userWeight, nil
}

func (s *UserStore) getLastWeight(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
) (string, error) {
	query := `
		SELECT weight FROM user_weight 
		WHERE user_id = $1 
		ORDER BY date DESC
		LIMIT 1
	`

	var w string
	err := tx.QueryRowContext(ctx, query, userID).Scan(&w)
	if err != nil {
		return "", err
	}

	return w, nil
}

// this return WITHOUT weight
func (s *UserStore) getUserAttr(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
) (UserAttributes, error) {
	query := `
		SELECT user_id, is_male, height, goal, weight_goal, age FROM user_attributes WHERE user_id = $1
	`
	var userAttr UserAttributes

	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&userAttr.UserID,
		&userAttr.IsMale,
		&userAttr.Height,
		&userAttr.Goal,
		&userAttr.WeightGoal,
		&userAttr.Age)
	if err != nil {
		return userAttr, err
	}
	return userAttr, nil
}

func (s *UserStore) updateUserWeight(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	weight float32,
) error {
	query := `
		UPDATE user_weight SET weight = $1 WHERE user_id = $2 AND date = CURRENT_DATE
	`

	_, err := tx.ExecContext(ctx, query, weight, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) addUserWeight(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	weight float32,
) error {
	query := `
		INSERT INTO user_weight (user_id, weight) VALUES ($1, $2)
	`

	_, err := tx.ExecContext(ctx, query, userID, weight)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

func (s *UserStore) addUserAttr(
	ctx context.Context,
	tx *sql.Tx,
	userAttr UserAttributes,
) error {
	query := `
    INSERT INTO user_attributes (user_id, is_male, height, goal, weight_goal, age) VALUES ($1, $2,$3,$4,$5,$6)
  `

	_, err := tx.ExecContext(
		ctx,
		query,
		userAttr.UserID,
		userAttr.IsMale,
		userAttr.Height,
		userAttr.Goal,
		userAttr.WeightGoal,
		userAttr.Age,
	)
	if err != nil {
		return err
	}
	return nil
}
