package main

import (
	"database/sql"
	"flag"
	"fmt"

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
		sqliteDatabase, _ := sql.Open("sqlite3", config.MealDatabasePath)
		defer sqliteDatabase.Close()
		allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
		strategy.MakeMealPlan(config, allMealsFromDatabase)
	} else {
		panic("No meal database path was provided in the configuration file.")
	}

	// database.GenerateDeterministicMealIDs()
}
