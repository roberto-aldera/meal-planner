package utilities

import (
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
)

func TestPrintMealDatabase(t *testing.T) {
	// First just check with empty database
	var emptyDatabase []database.Meal
	err := PrintMealDatabase(emptyDatabase)
	if err == nil {
		t.Fatal(err.Error())
	}

	allMealsFromDatabase := newDatabase(t)

	err = PrintMealDatabase(allMealsFromDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPrintExcludedMeals(t *testing.T) {
	// First just check with empty database
	var emptyMealMap map[int]database.Meal
	previousMealsToExclude := []int{103}
	err := PrintExcludedMeals(emptyMealMap, previousMealsToExclude)
	if err == nil {
		t.Fatal("Expected an error when using an empty meal map.")
	}

	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)

	// Check if no meals were excluded
	previousMealsToExclude = []int{}
	err = PrintExcludedMeals(mealMap, previousMealsToExclude)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check if excluded meal doesn't exist
	previousMealsToExclude = []int{1}
	err = PrintExcludedMeals(mealMap, previousMealsToExclude)
	if err == nil {
		t.Fatal("Expected an error when using a non-existent ID.")
	}

	// An finally check happy path runs as expected
	previousMealsToExclude = []int{103}
	err = PrintExcludedMeals(mealMap, previousMealsToExclude)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPrintMealPlan(t *testing.T) {
	// Check with a week plan full of empty meals
	var emptyMeal database.Meal
	weekPlan := []database.Meal{emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal}
	err := PrintMealPlan(weekPlan)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check with at least one real meal
	var realMeal database.Meal
	realMeal.ID = 123
	realMeal.MealName = "A tasty dish"

	weekPlan = []database.Meal{realMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal}
	err = PrintMealPlan(weekPlan)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check if week plan is incorrect length
	weekPlan = []database.Meal{emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal, emptyMeal}
	err = PrintMealPlan(weekPlan)
	if err == nil {
		t.Fatal("Expected an error when using a week meal plan of the incorrect length.")
	}
}
