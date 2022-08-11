package main

import (
	"flag"

	"github.com/roberto-aldera/meal-planner/strategy"
)

func main() {
	configFilePath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	strategy.MakeMealPlan(*configFilePath)
	// database.GenerateDeterministicMealIDs()
}
