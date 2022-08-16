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

func TestPrintMealDatabaseWithCategories(t *testing.T) {
	// First just check with empty database
	var emptyDatabase []database.Meal
	categories := []string{"Pasta"}
	err := PrintMealDatabaseWithCategories(emptyDatabase, categories)
	if err == nil {
		t.Fatal(err.Error())
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	// Try with empty category string
	var emptyCategories []string
	err = PrintMealDatabaseWithCategories(allMealsFromDatabase, emptyCategories)
	if err == nil {
		t.Fatal(err.Error())
	}

	categories = []string{"Pasta"}
	err = PrintMealDatabaseWithCategories(allMealsFromDatabase, categories)
	if err != nil {
		t.Fatal(err.Error())
	}
}
