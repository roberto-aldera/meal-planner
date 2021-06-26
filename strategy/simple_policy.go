package strategy

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/roberto-aldera/meal-planner/database"
)

func RunMe() {
	log.Println("Running policy...")
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()

	best_score := 100.0 // lower is better
	var best_meal_plan []database.Meal
	rand.Seed(time.Now().UTC().UnixNano())
	num_iterations := 100

	for i := 0; i < num_iterations; i++ {
		// Needed to load meals from database each time, as the meal picker deletes elements (can we pass a copy instead?)
		all_meals := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
		week_plan := pickRandomMeals(all_meals)
		meal_plan_score := calculateScore(week_plan)
		fmt.Println("Meal plan score:", meal_plan_score)
		if meal_plan_score < best_score {
			best_meal_plan = week_plan
			best_score = meal_plan_score
		}
	}
	fmt.Println("Best meal plan:")
	printMealPlan(best_meal_plan)
	fmt.Println("Score:", best_score)
}

func pickRandomMeals(all_meals []database.Meal) []database.Meal {
	// Pick 7 random meals for a start
	var week_plan []database.Meal // create empty plan
	initial_meal_idx := 0
	week_plan = append(week_plan, all_meals[rand.Intn(len(all_meals))])                 // initialise with a first dish
	all_meals = append(all_meals[:initial_meal_idx], all_meals[initial_meal_idx+1:]...) // erase dish from the possible options

	for len(week_plan) < 7 {
		idx := rand.Intn(len(all_meals))
		meal_under_test := all_meals[idx] // get a proposed meal
		week_plan = append(week_plan, meal_under_test)
		all_meals = append(all_meals[:idx], all_meals[idx+1:]...) // erase meal from available options
	}
	return week_plan
}

func printMealPlan(week_plan []database.Meal) {
	fmt.Println("Monday:   ", week_plan[0].Meal_name)
	fmt.Println("Tuesday:  ", week_plan[1].Meal_name)
	fmt.Println("Wednesday:", week_plan[2].Meal_name)
	fmt.Println("Thursday: ", week_plan[3].Meal_name)
	fmt.Println("Friday:   ", week_plan[4].Meal_name)
	fmt.Println("Saturday: ", week_plan[5].Meal_name)
	fmt.Println("Sunday:   ", week_plan[6].Meal_name)
}

// A big function to hold all the hand-written rules for now
// So tally things like cooking time, frequencies of dishes,
// complex things during the week, etc. and then score accordingly
func calculateScore(week_plan []database.Meal) float64 {
	// Higher numbers correspond to days where there is less time to cook
	time_penalties_per_day := [7]float64{1, 1, 10, 1, 10, -10, 1} // floats maybe?
	cooking_time_score := 0.0

	// Score for cooking times on days according to penalties
	for i := 0; i < len(week_plan); i++ {
		cooking_time_score += float64(week_plan[i].Cooking_time) * time_penalties_per_day[i]
	}

	return cooking_time_score
}