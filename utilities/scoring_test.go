package utilities

import (
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
)

func TestCalculateScore(t *testing.T) {
	filePath := "../default_config.json"
	config, _ := LoadConfiguration(filePath)
	config.PreferenceMealDaysOfWeek = []int{3}
	config.PreferenceMealIDs = []int{755}

	var emptyMeal database.Meal
	weekPlan := []database.Meal{emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal}

	_, err := CalculateScore(weekPlan, config)
	if err == nil {
		t.Fatal("Expected an error when using a week meal plan with at least one empty meal.")
	}

}
