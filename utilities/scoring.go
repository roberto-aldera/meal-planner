package utilities

import (
	"github.com/roberto-aldera/meal-planner/database"
)

// A big function to hold all the hand-written rules for now
// So tally things like cooking time, frequencies of dishes,
// complex things during the week, etc. and then score accordingly
func CalculateScore(weekPlan []database.Meal, config Config) float64 {
	// Higher numbers correspond to days where there is less time to cook
	cookingTimeScore := 0.0
	duplicateScore := 0.0
	finalMealPlanScore := 0.0
	lunchOnlyScore := 0.0

	// Score for cooking times on days according to penalties
	for i := 0; i < len(weekPlan); i++ {
		cookingTimeScore += float64(weekPlan[i].CookingTime) * config.DayWeights[i]
	}

	// Penalise duplicate categories within the same week
	tmpWeekPlan := make([]database.Meal, len(weekPlan))
	copy(tmpWeekPlan, weekPlan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmpWeekPlan); i++ {
		if visited[tmpWeekPlan[i].Category] {
			duplicateScore += config.DuplicatePenalty
		} else {
			visited[tmpWeekPlan[i].Category] = true
		}
	}

	// Penalise using lunch-only options for now
	for i := 0; i < len(weekPlan); i++ {
		if weekPlan[i].LunchOnly {
			lunchOnlyScore += config.LunchPenalty
		}
	}

	finalMealPlanScore = cookingTimeScore + duplicateScore + lunchOnlyScore
	return finalMealPlanScore
}
