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
	category := []string{"Pasta", "Soups", "Salad", "Healthy mix", "Curry", "Asian", "Meat with carb", "Rice/grains", "Pasta", "Pasta"}
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

	// Reduce default iterations when testing, no need for this to be so high
	config.NumberOfIterations = 50

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

	// Using a bad config
	config.PreferenceMealDaysOfWeek = []int{3}
	config.PreferenceMealIDs = []int{199}
	err = MakeMealPlan(config, newDatabase(t))
	if err != nil {
		fmt.Printf("MakeMealPlan failed: %s", err)
	}

	// Using a database with a duplicate meal
	config = newConfig(t)
	database_with_duplicate_meal := newDatabase(t)
	database_with_duplicate_meal[1].MealName = database_with_duplicate_meal[0].MealName

	err = MakeMealPlan(config, database_with_duplicate_meal)
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

func TestSetupMealMap(t *testing.T) {
	config := newConfig(t)
	_, err := setupMealMap(config, newDatabase(t))
	if err != nil {
		fmt.Printf("setupMealMap failed: %s", err)
	}
}

func TestAssignPreferences(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	_, err := assignPreferences(config, mealMap)
	if err != nil {
		fmt.Printf("assignPreferences failed: %s", err)
	}
}

func TestAssignPreferencesNonExistentMealIndex(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	config.PreferenceMealDaysOfWeek = []int{3}
	config.PreferenceMealIDs = []int{199}
	_, err := assignPreferences(config, mealMap)
	if err == nil {
		t.Fatal("Expected an error when deleting a non-existent meal.")
	}
}

func TestAssignPreferencesPrintExcludedMeals(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	config.PreviousMealsToExclude = []int{103}
	_, err := assignPreferences(config, mealMap)
	if err != nil {
		fmt.Printf("assignPreferences failed: %s", err)
	}

	// and try for non-existent meal index
	config.PreviousMealsToExclude = []int{199}
	_, err = assignPreferences(config, mealMap)
	if err == nil {
		t.Fatal("Expected an error when excluding a non-existent meal.")
	}
}

func TestAssignPreferencesRemoveNonExistentMeal(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	config.SpecialExclusions = []int{199}
	_, err := assignPreferences(config, mealMap)
	if err == nil {
		t.Fatal("Expected an error when excluding a non-existent meal.")
	}
}

func TestAssignPreferencesExcludeSoups(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	config.ExcludeSoups = true
	_, err := assignPreferences(config, mealMap)
	if err != nil {
		fmt.Printf("assignPreferences failed: %s", err)
	}

	// Try with non-soup category
	database_with_bad_soup_name := newDatabase(t)
	database_with_bad_soup_name[1].Category = "non-soup"
	mealMap, _ = setupMealMap(config, database_with_bad_soup_name)

	_, err = assignPreferences(config, mealMap)
	if err == nil {
		t.Fatal("Expected an error when excluding a non-existent meal.")
	}
}

func TestAssignPreferencesExcludeLunches(t *testing.T) {
	config := newConfig(t)
	mealMap, _ := setupMealMap(config, newDatabase(t))
	config.ExcludeLunches = true
	_, err := assignPreferences(config, mealMap)
	if err != nil {
		fmt.Printf("assignPreferences failed: %s", err)
	}

	// Try with no lunch meals
	database_with_no_lunch := newDatabase(t)
	database_with_no_lunch[9].LunchOnly = false
	mealMap, _ = setupMealMap(config, database_with_no_lunch)

	_, err = assignPreferences(config, mealMap)
	if err == nil {
		t.Fatal("Expected an error when excluding a non-existent meal.")
	}
}

func TestPickRandomMealsWithMapWithDuplicateMeals(t *testing.T) {
	config := newConfig(t)
	database_with_duplicate_meal := newDatabase(t)
	database_with_duplicate_meal[1].MealName = database_with_duplicate_meal[0].MealName
	mealMap, _ := setupMealMap(config, database_with_duplicate_meal)
	weekPlanWithRequests, _ := assignPreferences(config, mealMap)

	_, err := pickRandomMealsWithMap(mealMap, weekPlanWithRequests, config)
	if err != nil {
		fmt.Printf("pickRandomMealsWithMap failed: %s", err)
	}
}
