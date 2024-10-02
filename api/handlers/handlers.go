package handlers

import (
	"database/sql"

	"github.com/stanislavCasciuc/atom-fit/api/response"
)

type Handlers struct {
	resp *response.Responser
	db   *sql.DB
}

func New(resp response.Responser, db *sql.DB) *Handlers {
	return &Handlers{
		&resp,
		db,
	}
}
