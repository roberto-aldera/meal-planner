package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("meals.db")

	log.Println("Creating meals.db...")
	file, err := os.Create("meals.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("meals.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./meals.db")
	defer sqliteDatabase.Close()
	createTable(sqliteDatabase)
	insertMeal(sqliteDatabase, "Ragu", "Y")
	insertMeal(sqliteDatabase, "Amatriciana", "Y")
	insertMeal(sqliteDatabase, "Pasta al forno", "Y")
	displayMeals(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createMealTableSQL := `CREATE TABLE meals (
		"idMeal" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"tomato_based" TEXT		
	  );`

	log.Println("Create meals table...")
	statement, err := db.Prepare(createMealTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Meals table created")
}

func insertMeal(db *sql.DB, name string, tomato_based string) {
	log.Println("Inserting meal record ...")
	insertMealSQL := `INSERT INTO meals(name, tomato_based) VALUES (?, ?)`
	statement, err := db.Prepare(insertMealSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, tomato_based)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayMeals(db *sql.DB) {
	row, err := db.Query("SELECT * FROM meals ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var id int
		var name string
		var tomato_based string
		row.Scan(&id, &name, &tomato_based)
		log.Println("Meal: ", name, " ", tomato_based)
	}
}
