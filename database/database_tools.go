package database

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMe() {
	fmt.Println("Hello from database_tools!")

	// open database
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()

	displayEntries(sqliteDatabase)
	fmt.Println("All done!")
}

type Meal struct {
	Meal_name    string
	Cooking_time float32
	Category     string
	Lunch_only   bool
}

func countNumberOfRows(db *sql.DB) int {
	var num_rows int
	err := db.QueryRow("SELECT COUNT(*) FROM meals").Scan(&num_rows)
	switch {
	case err != nil:
		log.Fatal(err)
	default:
		// log.Printf("Number of rows are %d\n", num_rows)
	}
	return num_rows
}

func LoadDatabaseEntriesIntoContainer(db *sql.DB) []Meal {
	row, err := db.Query("SELECT Meal, Hours, Category, Lunch FROM meals ORDER BY Hours")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	num_rows := countNumberOfRows(db)
	var all_meals = make([]Meal, num_rows)
	i := 0
	for row.Next() {
		var meal Meal
		row.Scan(&meal.Meal_name, &meal.Cooking_time, &meal.Category, &meal.Lunch_only)
		all_meals[i] = meal
		i++
	}
	return all_meals
}

func displayEntries(db *sql.DB) {
	all_meals := LoadDatabaseEntriesIntoContainer(db)
	log.Println(all_meals)
}
