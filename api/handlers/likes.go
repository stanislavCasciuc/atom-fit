package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

var ErrLikeConflict = errors.New("user already liked this resource")

type ExerciseLikePayload struct {
	UserID     int64 `json:"user_id"`
	ExerciseID int64 `json:"exercise_id"`
}

// @LikeExercise	godoc
// @Summary		Like exercise
// @Description	Like exercise
// @Tags			likes
// @Accept			json
// @Produce		json
// @Param			exerciseID	path		int	true	"Exercise ID"
// @Success		200			{object}	response.SuccessResponse
// @Failure		400			{object}	error
// @Failure		500			{object}	error
// @Security		ApiKeyAuth
// @Router			/exercises/{exerciseID}/like [post]
func (h *Handlers) LikeExerciseHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
	idString := chi.URLParam(r, "exerciseID")
	exerciseID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
	}
	if err := h.store.ExercisesLikes.Create(r.Context(), u.ID, exerciseID); err != nil {
		switch err {
		case store.ErrConflict:
			h.resp.BadRequestError(w, r, ErrLikeConflict)
		default:
			h.resp.InternalServerError(w, r, err)
		}
		return
	}

	if err := response.WriteSuccess(w); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

// @UnLikeExercise	godoc
// @Summary		Unlike exercise
// @Description	Unlike exercise
// @Tags			likes
// @Accept			json
// @Produce		json
// @Param			exerciseID	path		int	true	"Exercise ID"
// @Success		200			{object}	response.SuccessResponse
// @Failure		400			{object}	error
// @Failure		500			{object}	error
// @Security		ApiKeyAuth
// @Router			/exercise/{exerciseID}/unlike [post]
func (h *Handlers) UnLikeExerciseHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
	idString := chi.URLParam(r, "exerciseID")
	exerciseID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
	}

	if err := h.store.ExercisesLikes.Delete(r.Context(), u.ID, exerciseID); err != nil {
		switch err {
		case store.ErrConflict:
			h.resp.BadRequestError(w, r, ErrLikeConflict)
		default:
			h.resp.InternalServerError(w, r, err)
		}
		return
	}

	if err := response.WriteSuccess(w); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
