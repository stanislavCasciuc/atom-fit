package handlers

import (
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type Handlers struct {
	resp   *response.Responser
	store  store.Storage
	config config.Config
}

func New(resp response.Responser, store store.Storage, config config.Config) *Handlers {
	return &Handlers{
		&resp,
		store,
		config,
	}
}
