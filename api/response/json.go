package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func (resp *Responser) ReadAndValidateJSON(w http.ResponseWriter, r *http.Request, data any) error {
	if err := ReadJSON(w, r, &data); err != nil {
		return err
	}
	if err := Validate.Struct(data); err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return WriteJSON(w, status, &envelope{Error: message})
}

type SuccessResponse struct {
	Status string `json:"status"`
}

func WriteSuccess(w http.ResponseWriter) error {
	return WriteJSON(w, http.StatusOK, SuccessResponse{Status: "success"})
}
