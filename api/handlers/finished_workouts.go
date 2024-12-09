package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type EndWorkoutPayload struct {
	WorkoutID int64  `json:"workout_id"`
	Duration  string `json:"duration"`
}

//	@EndWorkout		godoc
//	@Summary		End workout
//	@Description	End workout
//	@Tags			finished_workouts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		EndWorkoutPayload	true	"End workout payload"
//	@Success		200		{string}	string				"ok"
//	@Failure		400		{string}	string				"bad request"
//	@Failure		500		{string}	string				"internal server error"
//	@Security		ApiKeyAuth
//	@Router			/workouts/end [post]
func (h Handlers) EndWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	var payload EndWorkoutPayload
	if err := h.resp.ReadAndValidateJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	fn := store.FinishedWorkout{
		WorkoutID: payload.WorkoutID,
		Duration:  payload.Duration,
		UserID:    u.ID,
	}

	err := h.store.FinishedWorkouts.Create(r.Context(), &fn)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteSuccess(w); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
