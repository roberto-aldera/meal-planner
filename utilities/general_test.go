package utilities

import (
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
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

// Check that configuration loads as expected
func TestMakeMealMap(t *testing.T) {

	// First just check with empty database
	var emptyDatabase []database.Meal
	_, err := MakeMealMap(emptyDatabase)
	if err == nil {
		t.Fatal(err.Error())
	}

	allMealsFromDatabase := newDatabase(t)

	mealMap, err := MakeMealMap(allMealsFromDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(mealMap) != len(allMealsFromDatabase) {
		t.Fatal("Meal map does not contain all the meals stored in the original database.")
	}
}

func TestLoadMealRequestsAndUpdateMap(t *testing.T) {
	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	filePath := "../default_config.json"
	config, _ := LoadConfiguration(filePath)
	config.PreferenceMealDaysOfWeek = []int{3}
	config.PreferenceMealIDs = []int{103}
	lengthOfOriginalMap := len(mealMap)
	_, err := LoadMealRequestsAndUpdateMap(mealMap, config)
	if err != nil {
		t.Fatal(err)
	}

	if len(mealMap) >= lengthOfOriginalMap {
		t.Fatal("Updated meal map deletion did not occur as expected.")
	}
}

func TestRemoveSpecificMeals(t *testing.T) {
	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	lengthOfOriginalMap := len(mealMap)

	mealsToExclude := []int{1} // a non-existent ID
	err := RemoveSpecificMeals(mealMap, mealsToExclude)
	if err == nil {
		t.Fatal("Expected an error when deleting a non-existent meal.")
	}

	mealsToExclude = []int{103}
	err = RemoveSpecificMeals(mealMap, mealsToExclude)
	if err != nil {
		t.Fatal(err)
	}

	if len(mealMap) >= lengthOfOriginalMap {
		t.Fatal("Removing specific meals did not occur as expected.")
	}
}

func TestGetMealCategories(t *testing.T) {
	var emptyMealMap map[int]database.Meal
	_, err := GetMealCategories(emptyMealMap)

	if err == nil {
		t.Fatal("Expected an error when using an empty meal map.")
	}

	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	_, err = GetMealCategories(mealMap)

	if err != nil {
		t.Fatal("GetMealCategories has failed.")
	}
}

func TestGetMealsInCategory(t *testing.T) {
	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	_, err := GetMealsInCategory("Pet food", mealMap)

	if err == nil {
		t.Fatal("Expected an error when using a non-existent category.")
	}

	_, err = GetMealsInCategory("Pasta", mealMap)

	if err != nil {
		t.Fatal("GetMealsInCategory has failed.")
	}
}

func TestGetLunchMeals(t *testing.T) {
	var emptyMealMap map[int]database.Meal
	_, err := GetLunchMeals(emptyMealMap)

	if err == nil {
		t.Fatal("Expected an error when using an empty meal map.")
	}

	allMealsFromDatabase := newDatabase(t)
	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	_, err = GetLunchMeals(mealMap)

	if err != nil {
		t.Fatal("GetLunchMeals has failed.")
	}
}
