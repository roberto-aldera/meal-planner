package utilities

import (
	"database/sql"
	"testing"

	"github.com/roberto-aldera/meal-planner/database"
)

// Check that configuration loads as expected
func TestMakeMealMap(t *testing.T) {

	// First just check with empty database
	var emptyDatabase []database.Meal
	_, err := MakeMealMap(emptyDatabase)
	if err == nil {
		t.Fatal(err.Error())
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	mealMap, err := MakeMealMap(allMealsFromDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(mealMap) != len(allMealsFromDatabase) {
		t.Fatal("Meal map does not contain all the meals stored in the original database.")
	}
}

func TestLoadMealRequestsAndUpdateMap(t *testing.T) {
	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	filePath := "../default_config.json"
	config, _ := LoadConfiguration(filePath)
	config.PreferenceMealDaysOfWeek = []int{3}
	config.PreferenceMealIDs = []int{755}
	lengthOfOriginalMap := len(mealMap)
	_, err := LoadMealRequestsAndUpdateMap(mealMap, config)
	if err != nil {
		t.Fatal(err)
	}

	if len(mealMap) >= lengthOfOriginalMap {
		t.Fatal("Udpated meal map deletion did not occur as expected.")
	}
}

func TestRemoveSpecificMeals(t *testing.T) {
	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	mealMap, _ := MakeMealMap(allMealsFromDatabase)
	lengthOfOriginalMap := len(mealMap)

	mealsToExclude := []int{1} // a non-existent ID
	err := RemoveSpecificMeals(mealMap, mealsToExclude)
	if err == nil {
		t.Fatal("Expected an error when deleting a non-existent meal.")
	}

	mealsToExclude = []int{755}
	err = RemoveSpecificMeals(mealMap, mealsToExclude)
	if err != nil {
		t.Fatal(err)
	}

	if len(mealMap) >= lengthOfOriginalMap {
		t.Fatal("Removing specific meals did not occur as expected.")
	}
}
