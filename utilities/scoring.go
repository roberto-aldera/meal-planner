package utilities

import (
	"github.com/roberto-aldera/meal-planner/database"
)

// A big function to hold all the hand-written rules for now
// So tally things like cooking time, frequencies of dishes,
// complex things during the week, etc. and then score accordingly
func CalculateScore(weekPlan []database.Meal, config Config) float64 {
	mealPlanScore := 0.0

	// Score for meal complexity and cooking times for specified days
	for i := 0; i < len(weekPlan); i++ {
		// Meals that take longer to cook (are "complex") should be cooked on requested days
		// This penalises them being suggested when they weren't requested, and also penalises
		// them not being suggested when they were requested
		if config.ComplexMealRequested[i] && weekPlan[i].CookingTime < config.DefinitionOfLongMealPrepTimeHours {
			mealPlanScore += config.ScorePenalty
		} else if !config.ComplexMealRequested[i] && weekPlan[i].CookingTime > config.DefinitionOfLongMealPrepTimeHours {
			mealPlanScore += config.ScorePenalty
		}
		// Encourage "simple/quick" (user-defined in the database) meals on the requested days
		// And then also penalise days where a simple meal was not requested but was suggested
		if config.SimpleMealRequested[i] && !weekPlan[i].IsQuick {
			mealPlanScore += config.ScorePenalty
		} else if !config.SimpleMealRequested[i] && weekPlan[i].IsQuick {
			mealPlanScore += config.ScorePenalty
		}
	}

	// Penalise duplicate categories within the same week
	tmpWeekPlan := make([]database.Meal, len(weekPlan))
	copy(tmpWeekPlan, weekPlan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmpWeekPlan); i++ {
		if visited[tmpWeekPlan[i].Category] {
			mealPlanScore += config.ScorePenalty
		} else {
			visited[tmpWeekPlan[i].Category] = true
		}
	}
	return mealPlanScore
}
