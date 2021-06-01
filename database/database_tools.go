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

func displayEntries(db *sql.DB) {
	row, err := db.Query("SELECT Meal, Hours FROM meals ORDER BY Hours") // todo: using * for all wasn't working
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var Meal string
		var time string
		row.Scan(&Meal, &time)
		log.Println("Meal: ", Meal, " ", time)
	}
}
