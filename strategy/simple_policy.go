package strategy

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/roberto-aldera/meal-planner/database"
)

func RunMe() {
	log.Println("Running policy...")
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()
	all_meals := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
	pickRandomMeals(all_meals)
}

func pickRandomMeals(all_meals []database.Meal) {
	// Pick 7 random meals for a start
	var week_plan []database.Meal // create empty pasta plan
	initial_meal_idx := 0
	week_plan = append(week_plan, all_meals[0])                                         // initialise with a first dish
	all_meals = append(all_meals[:initial_meal_idx], all_meals[initial_meal_idx+1:]...) // erase dish from the possible options

	for len(week_plan) < 7 {
		idx := rand.Intn(len(all_meals))
		meal_under_test := all_meals[idx] // get a proposed meal
		week_plan = append(week_plan, meal_under_test)
		all_meals = append(all_meals[:idx], all_meals[idx+1:]...) // erase meal from available options
	}
	for _, meal := range week_plan {
		fmt.Println(meal.Meal_name)
	}
}
