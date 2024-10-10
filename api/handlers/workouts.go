package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

//	@CreateWorkout	godoc
//	@Summary		Create a new workout
//	@Description	Create a new workout
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	store.Workout	true	"Create Workout Payload"
//	@Success		204		"No Content"
//	@Security		ApiKeyAuth
//	@Router			/workouts [post]
func (h *Handlers) CreateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	var payload store.Workout
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}
	payload.UserID = u.ID

	err := h.store.Workouts.Create(r.Context(), &payload)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//	@GetAllWorkouts	godoc
//	@Summary		Get all workouts
//	@Description	Get all workouts
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.Workout
//	@Router			/workouts/ [get]
func (h *Handlers) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
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

	workouts, err := h.store.Workouts.GetAll(r.Context(), fq)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, workouts); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
