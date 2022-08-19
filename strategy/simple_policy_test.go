package strategy

import (
	"fmt"
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/utilities"
)

func newDatabase(t *testing.T) []database.Meal {
	identifiers := []int{101, 102, 103, 104, 105, 106, 107, 108, 109, 110}
	mealNames := []string{"Dish 1", "Dish 2", "Dish 3", "Dish 4", "Dish 5", "Dish 6", "Dish 7", "Dish 8", "Dish 9", "Dish 10"}
	cookingTime := []float64{1, 1, 0.5, 1, 1.5, 1.25, 1, 1.5, 0.75, 1}
	category := []string{"Pasta", "Soup", "Salad", "Healthy mix", "Curry", "Asian", "Meat with carb", "Rice/grains", "Pasta", "Pasta"}
	lunchOnly := []bool{false, false, false, false, false, false, false, false, false, true}
	isQuick := []bool{false, false, false, false, false, false, false, false, true, true}

	var mealDatabase []database.Meal

	for i := 0; i < len(identifiers); i++ {
		meal := database.Meal{
			ID:          identifiers[i],
			MealName:    mealNames[i],
			CookingTime: cookingTime[i],
			Category:    category[i],
			LunchOnly:   lunchOnly[i],
			IsQuick:     isQuick[i]}
		mealDatabase = append(mealDatabase, meal)
	}

	return mealDatabase
}

func newConfig(t *testing.T) (config utilities.Config) {
	configFilePath := "../default_config.json"

	config, err := utilities.LoadConfiguration(configFilePath)
	if err != nil {
		fmt.Printf("Configuration has failed to load: %s", err)
	}

	err = utilities.ValidateConfiguration(config)
	if err != nil {
		fmt.Printf("Configuration validation failed: %s", err)
	}
	return config
}
func TestMakeMealPlan(t *testing.T) {
	config := newConfig(t)
	err := MakeMealPlan(config, newDatabase(t))
	if err != nil {
		fmt.Printf("MakeMealPlan failed: %s", err)
	}
}

func TestMakeMealPlanWhenEmpty(t *testing.T) {
	config := newConfig(t)
	var emptyDatabase []database.Meal
	err := MakeMealPlan(config, emptyDatabase)
	if err == nil {
		t.Fatal("Expected an error when using an empty meal database.")
	}
}

func TestMakeMealPlanWithEmptyCategories(t *testing.T) {
	config := newConfig(t)
	database := newDatabase(t)

	for i := 0; i < len(database); i++ {
		database[i].Category = ""
	}

	err := MakeMealPlan(config, database)
	if err == nil {
		t.Fatal("Expected an error when using empty categories.")
	}
}
