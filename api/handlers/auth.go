package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

var exp = time.Duration(time.Hour * 24 * 3)

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
	if err := u.Password.Set(payload.Password); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := h.store.Users.CreateAndInvite(r.Context(), u, hashToken, exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			h.resp.BadRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			h.resp.BadRequestError(w, r, err)
		default:
			h.resp.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
