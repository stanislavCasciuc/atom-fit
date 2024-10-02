package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID        int64    `json:"id"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}
type password struct {
	Text *string
	Hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash
	p.Text = &text

	return nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, u.Email, u.Username, u.Password.Hash).
		Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) CreateAndInvite(
	ctx context.Context,
	user *User,
	token string,
	exp time.Duration,
) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		err := s.Create(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.createInvite(ctx, tx, user.ID, token, exp); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createInvite(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	token string,
	exp time.Duration,
) error {
	query := `
		INSERT INTO invitation (token, user_id, exp) 
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}
