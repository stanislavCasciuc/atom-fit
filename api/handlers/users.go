package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type registerUserPayload struct {
	Email    string `json:"email"    validation:"required,email"`
	Username string `json:"username" validation:"required,min=4,max=20"`
	Password string `json:"password" validation:"required,min=8"`
}

func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	u := &store.User{
		Email:    payload.Email,
		Username: payload.Username,
	}
	u.Password.Set(payload.Password)

	err := h.store.Users.Create(r.Context(), u)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
