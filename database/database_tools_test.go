package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createMealTable(db *sql.DB) {
	createMealTableSQL := `CREATE TABLE meals (
		"ID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Meal" TEXT,
		"Hours" integer,
		"Category" TEXT,
		"Lunch" integer,
		"Quick" integer		
	  );`

	fmt.Println("Create meals table...")
	statement, err := db.Prepare(createMealTableSQL)
	if err != nil {
		fmt.Print(err.Error())
	}
	statement.Exec()
	fmt.Println("Meals table created")
}

func insertMealIntoDatabase(db *sql.DB, name string, tomato_based string) {
	fmt.Println("Inserting meal record ...")
	insertMealSQL := `INSERT INTO meals(Meal, Category) VALUES (?, ?)`
	statement, err := db.Prepare(insertMealSQL)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(name, tomato_based)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestLoadDatabaseEntriesIntoContainer(t *testing.T) {
	os.Remove("meals.db")

	fmt.Println("Creating meals.db...")
	mealDatabasePath := t.TempDir() + "meals.db"
	file, err := os.Create(mealDatabasePath)
	if err != nil {
		fmt.Print(err.Error())
	}
	file.Close()
	fmt.Println("meals.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", mealDatabasePath)
	defer sqliteDatabase.Close()

	createMealTable(sqliteDatabase)
	insertMealIntoDatabase(sqliteDatabase, "Ragu", "Pasta")

	LoadDatabaseEntriesIntoContainer(sqliteDatabase)
}
