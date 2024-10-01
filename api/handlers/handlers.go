package handlers

import (
	"github.com/stanislavCasciuc/atom-fit/api/response"
)

type Handlers struct {
	resp *response.Responser
}

func New(resp response.Responser) *Handlers {
	return &Handlers{
		&resp,
	}
}
