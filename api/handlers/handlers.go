package handlers

import (
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/auth"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type Handlers struct {
	resp          response.Responser
	store         store.Storage
	config        config.Config
	authenticator auth.Authenticator
}

func New(
	resp response.Responser,
	store store.Storage,
	config config.Config,
	authenticator auth.Authenticator,
) *Handlers {
	return &Handlers{
		resp,
		store,
		config,
		authenticator,
	}
}
