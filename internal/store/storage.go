package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
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
	}
}

func New(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
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
