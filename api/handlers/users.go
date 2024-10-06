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

// ActivateUser godoc
//
//	@Summary		Activate a ActivateUser
//	@Description	Activate a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	ActivationPayload	true	"Activate User Payload"
//	@Success		204		"No Content"
//	@Router			/users/activate [post]
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

// @GetUserHandler	godoc
// @Summary		Get a user
// @Description	Get a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200	{object}	store.User
// @Security	 ApiKeyAuth
// @Router			/users [get]
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

// @GetUserWithAttrHandler	godoc
// @Summary				Get a user with attributes
// @Description			Get a user with attributes
// @Tags					users
// @Accept					json
// @Produce				json
// @Success				200	{object}	store.User
// @Security				ApiKeyAuth
// @Router					/users/attributes [get]
func (h *Handlers) GetUserWithAttrHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
	ua, err := h.store.Users.GetUserAttr(r.Context(), u.ID)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	u.UserAttr = *ua

	if err := response.WriteJSON(w, http.StatusOK, u); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

type LogWeightPayload struct {
	Weight float32 `json:"weight" validate:"required"`
}

// LogWeightHandler godoc
//
//	@Summary		Log a LogWeight
//	@Description	Log a user weight
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	LogWeightPayload	true	"Log Weight Payload"
//	@Success		204		"No Content"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/attributes/log/weight [post]
func (h *Handlers) LogWeightHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	var payload LogWeightPayload
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}
	if err := h.store.Users.AddUserWeight(r.Context(), u.ID, payload.Weight); err != nil {
		if err == store.ErrConflict {
			err = h.store.Users.UpdateUserWeight(r.Context(), u.ID, payload.Weight)
			if err != nil {
				h.resp.InternalServerError(w, r, err)
				return
			} else {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		h.resp.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
