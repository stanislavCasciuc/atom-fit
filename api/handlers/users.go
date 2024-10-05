package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/middleware"
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type ActivationPayload struct {
	Token string `json:"token"`
}

func (h *Handlers) ActivateUser(w http.ResponseWriter, r *http.Request) {
	var payload ActivationPayload
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	err := h.store.Users.Activate(r.Context(), payload.Token)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
	if err := response.WriteJSON(w, http.StatusOK, u); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

func (h *Handlers) GetUserFromCtx(r *http.Request) *store.User {
	ctx := r.Context()
	u, _ := ctx.Value(middleware.UserCtx).(*store.User)
	return u
}
