package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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
	ID          int
	MealName    string
	CookingTime float64
	Category    string
	LunchOnly   bool
}

func countNumberOfRows(db *sql.DB) int {
	var numRows int
	err := db.QueryRow("SELECT COUNT(*) FROM meals").Scan(&numRows)
	switch {
	case err != nil:
		log.Fatal(err)
	default:
		// log.Printf("Number of rows are %d\n", numRows)
	}
	return numRows
}

func LoadDatabaseEntriesIntoContainer(db *sql.DB) []Meal {
	row, err := db.Query("SELECT ID, Meal, Hours, Category, Lunch FROM meals ORDER BY Category")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	numRows := countNumberOfRows(db)
	var allMeals = make([]Meal, numRows)
	i := 0
	for row.Next() {
		var meal Meal
		row.Scan(&meal.ID, &meal.MealName, &meal.CookingTime, &meal.Category, &meal.LunchOnly)
		allMeals[i] = meal
		i++
	}
	return allMeals
}

func displayEntries(db *sql.DB) {
	allMeals := LoadDatabaseEntriesIntoContainer(db)
	log.Println(allMeals)
}

// Generate 3-digit unique IDs for each meal to be used to keep track of them in the database
func GenerateDeterministicMealIDs() {
	rand.Seed(42)
	numIDs := 899
	allIDs := rand.Perm(numIDs)
	numToPrint := 70
	for idx := range allIDs {
		allIDs[idx] += 100
	}
	for i := 0; i < numToPrint; i++ {
		fmt.Println(i, "->", allIDs[i])
	}
}
