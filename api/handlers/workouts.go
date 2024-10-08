package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
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
