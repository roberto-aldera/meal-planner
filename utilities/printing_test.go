package utilities

import (
	"database/sql"
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

	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	err = PrintMealDatabase(allMealsFromDatabase)
	if err != nil {
		t.Fatal(err.Error())
	}
}
