package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	ID        int64          `json:"id"`
	Email     string         `json:"email"`
	Username  string         `json:"username"`
	Password  password       `json:"-"`
	CreatedAt string         `json:"created_at"`
	IsActive  bool           `json:"is_active"`
	UserAttr  UserAttributes `json:"user_attr"`
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

		user.UserAttr.UserID = user.ID
		if err := s.addUserAttr(ctx, tx, user.UserAttr); err != nil {
			return err
		}

		if err := s.addUserWeight(ctx, tx, user.ID, user.UserAttr.Weight); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, username, password, created_at, is_active FROM users WHERE email = $1
	`

	var pass []byte
	u := &User{}
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Username, &pass, &u.CreatedAt, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.Password.Hash = pass
	return u, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, username, password, created_at, is_active FROM users WHERE id = $1
	`

	var pass []byte
	u := &User{}
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Email, &u.Username, &pass, &u.CreatedAt, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.Password.Hash = pass
	return u, nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		u, err := s.getFormInviteToken(ctx, tx, token)
		if err != nil {
			return err
		}

		u.IsActive = true
		if err := s.update(ctx, tx, u); err != nil {
			return err
		}

		err = s.deleteInvitation(ctx, tx, u.ID)
		if err != nil {
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

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := ` 
   DELETE FROM invitation WHERE user_id = $1
	`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getFormInviteToken(
	ctx context.Context,
	tx *sql.Tx,
	plainToken string,
) (*User, error) {
	query := `
		SELECT u.id, u.email, u.username,  u.is_active, u.created_at
		FROM users u
		JOIN invitation i ON u.id = i.user_id 
		WHERE i.token = $1
	`

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	u := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.IsActive,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return u, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, u *User) error {
	query := `
    UPDATE users SET email = $1, username = $2, is_active = $3 WHERE id = $4
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, u.Email, u.Username, u.IsActive, u.ID)
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
