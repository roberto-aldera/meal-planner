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
