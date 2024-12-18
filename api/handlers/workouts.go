package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/mailer/pagination"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type CreateWorkoutPayload struct {
	Name             string                         `json:"name"          validate:"required"`
	Description      string                         `json:"description"`
	TutorialLink     string                         `json:"tutorial_link"`
	ExercisesWorkout []CreateWorkoutExercisePayload `json:"exercises"     validate:"required"`
}
type CreateWorkoutExercisePayload struct {
	ExerciseID int64 `json:"exercise_id" validate:"required"`
	Duration   int   `json:"duration"    validate:"required"`
}

//	@CreateWorkout	godoc
//	@Summary		Create a new workout
//	@Description	Create a new workout
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	CreateWorkoutPayload	true	"Create Workout Payload"
//	@Success		204		"No Content"
//	@Security		ApiKeyAuth
//	@Router			/workouts [post]
func (h *Handlers) CreateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	var payload CreateWorkoutPayload
	if err := response.ReadJSON(w, r, &payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(payload); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}
	workout := store.Workout{
		UserID:       u.ID,
		Name:         payload.Name,
		Description:  payload.Description,
		TutorialLink: payload.TutorialLink,
	}
	for _, e := range payload.ExercisesWorkout {
		workout.WorkoutExercises = append(workout.WorkoutExercises, store.WorkoutExercises{
			ExerciseID: e.ExerciseID,
			Duration:   e.Duration,
		})
	}
	err := h.store.Workouts.Create(r.Context(), &workout)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteSuccess(w); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

type WorkoutsResponse struct {
	Data       []store.Workout `json:"data"`
	TotalCount int             `json:"total_count"`
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
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	WorkoutsResponse
//	@Security		ApiKeyAuth
//	@Router			/workouts/ [get]
func (h *Handlers) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)
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

	workouts, totalCount, err := h.store.Workouts.GetAll(r.Context(), fq, u.ID)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	resp := WorkoutsResponse{
		Data:       workouts,
		TotalCount: totalCount,
	}
	if err := response.WriteJSON(w, http.StatusOK, resp); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

//	@GetUserWorkouts	godoc
//	@Summary			Get user workouts
//	@Description		Get user workouts
//	@Tags				workouts
//	@Accept				json
//	@Produce			json
//	@Param				limit	query		int		false	"Limit"
//	@Param				offset	query		int		false	"Offset"
//	@Param				sort	query		string	false	"Sort"
//	@Param				userID	path		int		false	"userID"
//	@Success			200		{object}	[]store.Workout
//	@Security			ApiKeyAuth
//	@Router				/workouts/user/{userID} [get]
func (h *Handlers) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}
	fq := pagination.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err = fq.Parse(r)
	if err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	if err := response.Validate.Struct(fq); err != nil {
		h.resp.BadRequestError(w, r, err)
		return
	}

	workouts, err := h.store.Workouts.GetUsersWorkouts(r.Context(), fq, userID)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, workouts); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

//	@GetWorkout		godoc
//	@Summary		Get workout by ID
//	@Description	Get workout by ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	path		int	true	"Workout ID"
//	@Success		200			{object}	store.Workout
//	@Router			/workouts/{workoutID} [get]
func (h *Handlers) GetWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	workoutID := chi.URLParam(r, "workoutID")
	if workoutID == "" {
		h.resp.BadRequestError(w, r, errors.New("missing workoutID"))
		return
	}
	id, err := strconv.ParseInt(workoutID, 10, 64)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	workout, err := h.store.Workouts.GetByID(ctx, id)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
	if workout == nil {
		h.resp.NotFoundErorr(w, r, errors.New("resource not found"))
		return
	}

	we, err := h.store.Workouts.GetWorkoutExercises(ctx, id)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	workout.WorkoutExercises = we

	if err := response.WriteJSON(w, http.StatusOK, workout); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}

//	@DeleteWorkout	godoc
//	@Summary		Delete workout by ID
//	@Description	Delete workout by ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	path	int	true	"Workout ID"
//	@Success		204			"No Content"
//	@Security		ApiKeyAuth
//	@Router			/workouts/{workoutID} [delete]
func (h *Handlers) DeleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
}

//	@PatchWorkout	godoc
//	@Summary		Update workout by ID
//	@Description	Update workout by ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	path	int	true	"Workout ID"
//	@Success		204			"No Content"
//	@Security		ApiKeyAuth
//	@Router			/workouts/{workoutID} [patch]
func (h *Handlers) PatchWorkoutHandler(w http.ResponseWriter, r *http.Request) {
}
