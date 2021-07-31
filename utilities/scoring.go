package utilities

import (
	"github.com/roberto-aldera/meal-planner/database"
)

// A big function to hold all the hand-written rules for now
// So tally things like cooking time, frequencies of dishes,
// complex things during the week, etc. and then score accordingly
func CalculateScore(week_plan []database.Meal, config Config) float64 {
	// Higher numbers correspond to days where there is less time to cook
	cooking_time_score := 0.0
	duplicate_score := 0.0
	final_meal_plan_score := 0.0
	lunch_only_score := 0.0

	// Score for cooking times on days according to penalties
	for i := 0; i < len(week_plan); i++ {
		cooking_time_score += float64(week_plan[i].Cooking_time) * config.Day_weights[i]
	}

	// Penalise duplicate categories within the same week
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Category] {
			duplicate_score += config.Duplicate_penalty
		} else {
			visited[tmp_week_plan[i].Category] = true
		}
	}

	// Penalise using lunch-only options for now
	for i := 0; i < len(week_plan); i++ {
		if week_plan[i].Lunch_only {
			lunch_only_score += config.Lunch_penalty
		}
	}

	final_meal_plan_score = cooking_time_score + duplicate_score + lunch_only_score
	return final_meal_plan_score
}
