package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
)

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "OK"}
	if err := response.WriteJSON(w, http.StatusOK, data); err != nil {
		response.WriteJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
