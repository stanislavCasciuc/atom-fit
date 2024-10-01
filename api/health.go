package api

import "net/http"

func (a *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "OK"}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
