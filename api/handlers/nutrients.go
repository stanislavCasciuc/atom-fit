package handlers

import (
	"net/http"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/nutrients"
)

//	@GetMacronutrientsGoalPerDayHandler	godoc
//	@Summary							Get macronutrients goal per day
//	@Description						Get macronutrients goal per day for the user
//	@Tags								nutrients
//	@Accept								json
//	@Produce							json
//	@Security							ApiKeyAuth
//	@Success							200	{object}	nutrients.UserNutrientsGoal
//	@Router								/nutrients/daily-goal [get]
func (h *Handlers) GetMacronutrientsGoalPerDayHandler(w http.ResponseWriter, r *http.Request) {
	u := h.GetUserFromCtx(r)

	userAttr, err := h.store.Users.GetUserAttr(r.Context(), u.ID)
	if err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}

	goal := nutrients.CalculateMacronutrients(*userAttr)

	if err := response.WriteJSON(w, http.StatusOK, goal); err != nil {
		h.resp.InternalServerError(w, r, err)
		return
	}
}
