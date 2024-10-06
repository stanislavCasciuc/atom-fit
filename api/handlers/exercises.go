package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

//	@CreateExerciseHandler	godoc
//	@Summary				Create a new Exercise
//	@Description			Create a new Exercise
//	@Tags					exercises
//	@Accept					json
//	@Produce				json
//	@Param					payload	body		store.Exercise	true	"Exercise Payload"
//	@Success				201		{object}	store.Exercise
//	@Security				ApiKeyAuth
//	@Router					/exercises [post]
func (h *Handlers) CreateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	var payload store.Exercise
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}
	payload.UserID = u.ID

	if err := h.store.Exercises.Create(r.Context(), &payload); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusCreated, payload); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

//	@GetExercisesHandler	godoc
//	@Summary				Get Exercise by id from param
//	@Description			Get Exercise by id from param
//	@Tags					exercises
//	@Accept					json
//	@Produce				json
//	@Param					exerciseID	path		int	true	"Exercise ID"
//	@Success				200			{object}	store.Exercise
//	@Router					/exercises/{exerciseID} [get]
func (h *Handlers) GetExerciseHandler(w http.ResponseWriter, r *http.Request) {
	exerciseID := chi.URLParam(r, "exerciseID")
	id, err := strconv.ParseInt(exerciseID, 10, 64)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	e, err := h.store.Exercises.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			h.resp.BadRequestError(w, r, err)
		default:
			h.resp.InternalServerError(w, r, err)
		}
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, e); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
