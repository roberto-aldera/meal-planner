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
	all_meals_from_database := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	best_score := 100.0 // lower is better
	var best_meal_plan []database.Meal
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())
	num_iterations := 10000

	for i := 0; i < num_iterations; i++ {
		// Need to make copy of slice, as modifications affect underlying array
		tmp_all_meals := make([]database.Meal, len(all_meals_from_database))
		copy(tmp_all_meals, all_meals_from_database)

		week_plan := pickRandomMeals(tmp_all_meals)
		meal_plan_score := calculateScore(week_plan)
		if meal_plan_score < best_score {
			best_meal_plan = week_plan
			best_score = meal_plan_score
		}
	}

	fmt.Println("Best meal plan after", num_iterations, "iterations from a total of", len(all_meals_from_database), "meals:")
	printMealPlan(best_meal_plan)
	fmt.Println("Score:", best_score)
}

func pickRandomMeals(all_meals []database.Meal) []database.Meal {
	// Pick 7 random meals for a start
	var week_plan []database.Meal // create empty plan
	initial_meal_idx := rand.Intn(len(all_meals))
	week_plan = append(week_plan, all_meals[initial_meal_idx])                          // initialise with a first dish
	all_meals = append(all_meals[:initial_meal_idx], all_meals[initial_meal_idx+1:]...) // erase dish from the possible options

	for len(week_plan) < 7 {
		idx := rand.Intn(len(all_meals))
		meal_under_test := all_meals[idx] // get a proposed meal
		week_plan = append(week_plan, meal_under_test)
		all_meals = append(all_meals[:idx], all_meals[idx+1:]...) // erase meal from available options
	}

	// Debug: check for duplicates
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Meal_name] {
			fmt.Println("*** Dupilcate found:", tmp_week_plan[i].Meal_name)
		} else {
			visited[tmp_week_plan[i].Meal_name] = true
		}
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
	time_penalties_per_day := [7]float64{1, 1, 30, 1, 30, -10, 5}
	cooking_time_score := 0.0
	duplicate_score := 0.0
	final_meal_plan_score := 0.0
	lunch_only_score := 0.0

	// Score for cooking times on days according to penalties
	for i := 0; i < len(week_plan); i++ {
		cooking_time_score += float64(week_plan[i].Cooking_time) * time_penalties_per_day[i]
	}

	// Penalise duplicate categories within the same week
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Category] {
			duplicate_score += 10
		} else {
			visited[tmp_week_plan[i].Category] = true
		}
	}

	// Penalise using lunch-only options for now
	for i := 0; i < len(week_plan); i++ {
		if week_plan[i].Lunch_only {
			lunch_only_score += 100
		}
	}

	final_meal_plan_score = cooking_time_score + duplicate_score + lunch_only_score
	return final_meal_plan_score
}
