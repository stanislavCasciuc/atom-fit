package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type WorkoutReviewPayload struct {
	Title   string `json:"title"   validate:"required,max=255"`
	Rating  int    `json:"rating"  validate:"required,oneof=1 2 3 4 5"`
	Content string `json:"content" validate:"required"`
}

// @ReviewWorkout	godoc
// @Summary		Review workout
// @Description	Review workout
// @Tags			reviews
// @Accept			json
// @Produce		json
// @Param			workoutID	path		int						true	"Workout ID"
// @Param			payload		body		WorkoutReviewPayload	true	"Review workout payload"
// @Success		200			{object}	store.WorkoutReview
// @Failure		400			{object}	error
// @Security		ApiKeyAuth
// @Router			/reviews/workout/{workoutID} [post]
func (h *Handlers) ReviewWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
	idString := chi.URLParam(r, "workoutID")
	workoutID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
	var payload WorkoutReviewPayload
	if err := h.resp.ReadAndValidateJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	wr := &store.WorkoutReview{
		WorkoutID: workoutID,
		Title:     payload.Title,
		Content:   payload.Content,
		Rating:    payload.Rating,
		UserID:    u.ID,
	}
	err = h.store.Reviews.CreateWorkout(r.Context(), wr)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, wr); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

// @GetWorkoutReviews	godoc
// @Summary			Get workout reviews
// @Description		Get workout reviews
// @Tags				reviews
// @Accept				json
// @Produce			json
// @Param				workoutID	path		int	true	"Workout ID"
// @Success			200			{object}	[]store.WorkoutReviewWithMetadata
// @Failure			400			{object}	error
// @Router				/reviews/workout/{workoutID} [get]
func (h *Handlers) GetWorkoutReviewsHandler(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "workoutID")
	workoutID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	wr, err := h.store.Reviews.Get(r.Context(), workoutID)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, wr); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
