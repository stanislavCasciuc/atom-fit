package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

// @GetAllExercises	godoc
// @Summary			Get all Exercises
// @Description		Get all Exercises
// @Tags				exercises
// @Accept				json
// @Produce			json
//
// @Param				since	query		string	false	"Since"
// @Param				until	query		string	false	"Until"
// @Param				limit	query		int		false	"Limit"
// @Param				offset	query		int		false	"Offset"
// @Param				sort	query		string	false	"Sort"
// @Param				tags	query		string	false	"Tags"
// @Param				search	query		string	false	"Search"
//
// @Success			200		{object}	[]store.Exercise
// @Router				/exercises [get]
func (h *Handlers) GetAllExercisesHandler(w http.ResponseWriter, r *http.Request) {
	fq := pagination.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(fq); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	log.Println(fq)

	exercises, err := h.store.Exercises.GetAll(r.Context(), fq)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, exercises); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

// @CreateExerciseHandler	godoc
// @Summary				Create a new Exercise
// @Description			Create a new Exercise
// @Tags					exercises
// @Accept					json
// @Produce				json
// @Param					payload	body		store.Exercise	true	"Exercise Payload"
// @Success				201		{object}	store.Exercise
// @Security				ApiKeyAuth
// @Router					/exercises [post]
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

// @GetExercisesHandler	godoc
// @Summary				Get Exercise by id from param
// @Description			Get Exercise by id from param
// @Tags					exercises
// @Accept					json
// @Produce				json
// @Param					exerciseID	path		int	true	"Exercise ID"
// @Success				200			{object}	store.Exercise
// @Router					/exercises/{exerciseID} [get]
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