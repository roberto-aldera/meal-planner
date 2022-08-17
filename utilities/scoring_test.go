package utilities

import (
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
)

func TestCalculateScore(t *testing.T) {
	filePath := "../default_config.json"
	config, _ := LoadConfiguration(filePath)
	// config.PreferenceMealDaysOfWeek = []int{3}
	// config.PreferenceMealIDs = []int{755}

	var emptyMeal database.Meal
	weekPlan := []database.Meal{emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal}

	_, err := CalculateScore(weekPlan, config)
	if err == nil {
		t.Fatal("Expected an error when using a week meal plan with at least one empty meal.")
	}

	var complexMeal database.Meal
	complexMeal.ID = 123
	complexMeal.MealName = "Complex dish"
	complexMeal.CookingTime = 3

	var standardMeal database.Meal
	standardMeal.ID = 124
	standardMeal.MealName = "Standard dish"
	standardMeal.CookingTime = 1

	var quickMeal database.Meal
	quickMeal.ID = 125
	quickMeal.MealName = "Quick dish"
	quickMeal.CookingTime = 0.25
	quickMeal.IsQuick = true

	config.ComplexMealRequested = []bool{true, false, false, false, false, false, false}
	config.SimpleMealRequested = []bool{false, false, false, false, false, false, true}

	weekPlan = []database.Meal{standardMeal, complexMeal, quickMeal, complexMeal, complexMeal, complexMeal, complexMeal}

	_, err = CalculateScore(weekPlan, config)
	if err != nil {
		t.Fatal(err.Error())
	}
}
