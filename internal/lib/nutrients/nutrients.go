package nutrients

import "github.com/stanislavCasciuc/atom-fit/internal/store"

type UserNutrientsGoal struct {
	Calories     float32
	Proteins     float32
	Fats         float32
	Carbohydrats float32
}

func CalculateMacronutrients(userAttr store.UserAttributes) UserNutrientsGoal {
	// Calculate BMR using Mifflin-St Jeor Equation
	var BMR float32
	if userAttr.IsMale {
		BMR = 10*userAttr.Weight + 6.25*float32(userAttr.Height) - 5*float32(userAttr.Age) + 5
	} else {
		BMR = 10*userAttr.Weight + 6.25*float32(userAttr.Height) - 5*float32(userAttr.Age) - 161
	}

	// Assuming a moderate activity level (TDEE multiplier ~1.55)
	TDEE := BMR * 1.55

	// Adjust TDEE based on user's goal
	switch userAttr.Goal {
	case "lose":
		TDEE -= 500 // To lose weight, we aim for a 500-calorie deficit.
	case "gain":
		TDEE += 500 // To gain weight, we aim for a 500-calorie surplus.
	case "maintain":
		// No change for maintenance.
	}

	// Set protein intake (grams per kg of body weight)
	var proteinPerKg float32 = 2.0 // Setting a baseline protein intake.
	if userAttr.Goal == "gain" {
		proteinPerKg = 2.2
	} else if userAttr.Goal == "lose" {
		proteinPerKg = 1.8
	}
	proteinGrams := proteinPerKg * userAttr.Weight

	// Calculate protein calories (1g of protein = 4 calories)
	proteinCalories := proteinGrams * 4

	// Set fat intake (~30% of total TDEE)
	fatCalories := TDEE * 0.30
	fatGrams := fatCalories / 9 // 1g of fat = 9 calories

	// Remaining calories for carbohydrates
	remainingCalories := TDEE - (proteinCalories + fatCalories)
	carbGrams := remainingCalories / 4 // 1g of carbohydrates = 4 calories

	// Return the calculated macronutrients
	return UserNutrientsGoal{
		Calories:     TDEE,
		Proteins:     proteinGrams,
		Fats:         fatGrams,
		Carbohydrats: carbGrams,
	}
}
