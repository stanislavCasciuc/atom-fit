package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
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
