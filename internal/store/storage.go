package store

import (
	"context"
	"database/sql"
	"time"
)

var QueryTimeDuration = time.Second * 5

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
	}
}

func New(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
	}
}
