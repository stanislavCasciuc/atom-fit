package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

var exp = time.Duration(time.Hour * 24 * 3)

type registerUserPayload struct {
	Email    string `json:"email"    validation:"required,email"`
	Username string `json:"username" validation:"required,min=4,max=20"`
	Password string `json:"password" validation:"required,min=8"`
}

type TokenResponse struct {
	Token string `json:"token"`
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

	res := TokenResponse{
		Token: plainToken,
	}

	go mailer.SendVerifyUser(u.Username, h.config.Mail.Addr, plainToken, h.config.Mail)

	if err := response.WriteJSON(w, http.StatusOK, res); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

type LoginPayload struct {
	Email    string `json:"email"    validation:"required,email"`
	Password string `json:"password" validation:"required,min=8"`
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	u, err := h.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			h.resp.NotFoundErorr(w, r, err)
		default:
			h.resp.InternalServerError(w, r, err)
		}
		return
	}

	claims := jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(h.config.Auth.Iat).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Aud,
		"aud": h.config.Auth.Aud,
	}

	token, err := h.authenticator.GenerateToken(claims)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, token); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
