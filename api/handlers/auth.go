package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

var exp = time.Duration(time.Hour * 24 * 3)

type registerUserPayload struct {
	Email      string  `json:"email"       validate:"required,email"`
	Username   string  `json:"username"    validate:"required,min=4,max=20"`
	Password   string  `json:"password"    validate:"required,min=8"`
	IsMale     bool    `json:"is_male"     validate:"required"`
	Height     int     `json:"height"      validate:"max=250,min=100"`
	Goal       string  `json:"goal"        validate:"required,oneof=lose gain maintain"`
	WeightGoal float32 `json:"weight_goal"`
	Weight     float32 `json:"weight"`
	Age        int     `json:"age"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// RegisterUserHandler godoc
//
//	@Summary		Register a new
//	@Description	Register a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"Register User Payload"
//	@Success		200		{object}	TokenResponse
//	@Router			/auth/register [post]
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
		UserAttr: store.UserAttributes{
			IsMale:     payload.IsMale,
			Height:     payload.Height,
			Goal:       payload.Goal,
			WeightGoal: payload.WeightGoal,
			Weight:     payload.Weight,
			Age:        payload.Age,
		},
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

	// go mailer.SendVerifyUser(u.Username, h.config.Mail.Addr, plainToken, h.config.Mail)

	if err := response.WriteJSON(w, http.StatusOK, res); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

type LoginPayload struct {
	Email    string `json:"email"    validation:"required,email"`
	Password string `json:"password" validation:"required,min=8"`
}

// LoginHandler godoc
//
//	@Summary		LoginHandler
//	@Description	LoginHandler
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginPayload	true	"Login Payload"
//	@Success		200		{object}	TokenResponse
//	@Router			/auth/login [post]
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
	tokenStruct := TokenResponse{
		Token: token,
	}
	if err := response.WriteJSON(w, http.StatusOK, tokenStruct); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
