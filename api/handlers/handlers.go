package handlers

import (
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type Handlers struct {
	resp  *response.Responser
	store store.Storage
}

func New(resp response.Responser, store store.Storage) *Handlers {
	return &Handlers{
		&resp,
		store,
	}
}
