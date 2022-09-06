package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/strategy"
	"github.com/roberto-aldera/meal-planner/utilities"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	config, err := utilities.LoadConfiguration(*configFilePath)
	if err != nil {
		fmt.Printf("Configuration has failed to load: %s", err)
	}

	err = utilities.ValidateConfiguration(config)
	if err != nil {
		fmt.Printf("Configuration validation failed: %s", err)
	}

	// Load meals from database and print out all candidates
	if config.MealDatabasePath != "" {
		sqliteDatabase, err := sql.Open("sqlite3", config.MealDatabasePath)
		if err != nil {
			fmt.Printf("Failure in opening meal database: %s", err)
		}
		fmt.Printf("sqliteDatabase.Stats(): %v\n", sqliteDatabase.Stats())
		defer sqliteDatabase.Close()
		allMealsFromDatabase, err := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
		if err != nil {
			fmt.Printf("Failure in loading meals from database: %s", err)
		}
		strategy.MakeMealPlan(config, allMealsFromDatabase)
	} else {
		panic("No meal database path was provided in the configuration file.")
	}

	// database.GenerateDeterministicMealIDs()
}
