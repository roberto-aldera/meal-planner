package strategy

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/utilities"
)

type specific_meal struct {
	Meal_ID_idx int
	Day_of_week int
}

func RunMe() {
	log.Println("Running policy...")

	// Load meals from database and print out all candidates
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()
	all_meals_from_database := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
	utilities.PrintMealDatabase(all_meals_from_database)

	// Build config
	var config utilities.Config
	config.Number_of_iterations = 1 //00000
	config.Day_weights = [7]float64{1, 1, 30, 1, 30, -10, 30}
	config.Minimum_score = 10000
	config.Duplicate_penalty = 100
	config.Lunch_penalty = 100

	// Handle pre-selected meals
	meal_IDs := []int{197, 752, 255}
	var meal_ID_idx []int
	meal_days_of_the_week := []int{0, 1, 4}
	var meals_to_load []specific_meal
	// Quick check that the inputs are legal
	if len(meal_IDs) == len(meal_days_of_the_week) {
		// Do a conversion from meal_ID to index
		// Need to go through all meal IDs, and find IDs that are in both lists.
		// Then save the indices where these occured, and proceed.
		for _, id := range meal_IDs {
			for idx, meal := range all_meals_from_database {
				if meal.ID == id {
					meal_ID_idx = append(meal_ID_idx, idx)
				}
			}
		}

		for idx := range meal_IDs {
			var meal_to_load specific_meal
			meal_to_load.Meal_ID_idx, meal_to_load.Day_of_week = meal_ID_idx[idx], meal_days_of_the_week[idx] //meal_IDs[idx], meal_days_of_the_week[idx]
			meals_to_load = append(meals_to_load, meal_to_load)
		}
	}
	fmt.Println("Your requested meals:")
	for _, meal := range meals_to_load {
		fmt.Println("Day of the week:", meal.Day_of_week, "- meal:", all_meals_from_database[meal.Meal_ID_idx].Meal_name)
	}

	best_score := config.Minimum_score // lower is better
	var best_meal_plan []database.Meal
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < config.Number_of_iterations; i++ {
		// Need to make copy of slice, as modifications affect underlying array
		tmp_all_meals := make([]database.Meal, len(all_meals_from_database))
		copy(tmp_all_meals, all_meals_from_database)

		week_plan := pickRandomMeals(tmp_all_meals, meals_to_load, config)
		meal_plan_score := calculateScore(week_plan, config)
		if meal_plan_score < best_score {
			best_meal_plan = week_plan
			best_score = meal_plan_score
		}
	}

	if len(best_meal_plan) == 7 {
		fmt.Println("Best meal plan after", config.Number_of_iterations, "iterations from a total of", len(all_meals_from_database), "meals:")
		utilities.PrintMealPlan(best_meal_plan)
		fmt.Println("Score:", best_score)
	} else {
		fmt.Println("No valid meal plan was possible with the provided requirements.")
	}
}

func RunMeWithMap() {
	log.Println("Running policy...")

	// Load meals from database and print out all candidates
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()
	all_meals_from_database := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)
	utilities.PrintMealDatabase(all_meals_from_database)

	meal_map := make_meal_map(all_meals_from_database)

	// Build config
	var config utilities.Config
	config.Number_of_iterations = 100000
	config.Day_weights = [7]float64{1, 1, 30, 1, 30, -10, 30}
	config.Minimum_score = 10000
	config.Duplicate_penalty = 100
	config.Lunch_penalty = 100

	best_score := config.Minimum_score // lower is better
	var best_meal_plan []database.Meal
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < config.Number_of_iterations; i++ {
		week_plan := pickRandomMealsWithMap(meal_map, config)
		meal_plan_score := calculateScore(week_plan, config)
		if meal_plan_score < best_score {
			best_meal_plan = week_plan
			best_score = meal_plan_score
		}
	}

	if len(best_meal_plan) == 7 {
		fmt.Println("Best meal plan after", config.Number_of_iterations, "iterations from a total of", len(all_meals_from_database), "meals:")
		utilities.PrintMealPlan(best_meal_plan)
		fmt.Println("Score:", best_score)
	} else {
		fmt.Println("No valid meal plan was possible with the provided requirements.")
	}
}

func get_next_empty_slot(week_plan []database.Meal) int {
	for idx, item := range week_plan {
		if item.Meal_name == "" {
			return idx
		}
	}
	return -1
}

func pickRandomMeals(all_meals []database.Meal, meals_to_load []specific_meal, config utilities.Config) []database.Meal {
	week_plan := make([]database.Meal, 7)

	// Load pre-selected meals into meal plan
	for _, meal_to_load := range meals_to_load {
		week_plan[meal_to_load.Day_of_week] = all_meals[meal_to_load.Meal_ID_idx]
	}

	// Erase dishes from the possible options - only once all have been added as this alters all_meals
	// This is broken, the indices used to arrange change with each deletion
	// Might be better to do this with a map
	for _, meal_to_load := range meals_to_load {
		all_meals = append(all_meals[:meal_to_load.Meal_ID_idx], all_meals[meal_to_load.Meal_ID_idx+1:]...)
	}

	next_idx := get_next_empty_slot(week_plan)
	for next_idx >= 0 {
		idx := rand.Intn(len(all_meals))
		meal_under_test := all_meals[idx] // get a proposed meal
		week_plan[next_idx] = meal_under_test
		all_meals = append(all_meals[:idx], all_meals[idx+1:]...) // erase meal from available options
		next_idx = get_next_empty_slot(week_plan)
	}

	// Debug: check for duplicates
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Meal_name] {
			fmt.Println("*** Duplicate found:", tmp_week_plan[i].Meal_name)
		} else {
			visited[tmp_week_plan[i].Meal_name] = true
		}
	}

	return week_plan
}

func pickRandomMealsWithMap(meal_map map[int]database.Meal, config utilities.Config) []database.Meal {
	week_plan := make([]database.Meal, 7)

	// Store map keys in a slice, and get N random items from this slice to use in the plan (to avoid picking duplicates)
	// Later: remove elements from slice that are keys already hand-picked by user
	slice_of_keys := make([]int, 0)
	for key := range meal_map {
		slice_of_keys = append(slice_of_keys, key)
	}

	random_indices := rand.Perm(len(meal_map))
	key_subset := make([]int, 0)
	for i := 0; i < len(week_plan); i++ {
		key_subset = append(key_subset, slice_of_keys[random_indices[i]])
	}

	for idx := 0; idx < len(week_plan); idx++ {
		meal_under_test := meal_map[key_subset[idx]] // get a proposed meal
		week_plan[idx] = meal_under_test
	}

	// Debug: check for duplicates
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Meal_name] {
			fmt.Println("*** Duplicate found:", tmp_week_plan[i].Meal_name)
		} else {
			visited[tmp_week_plan[i].Meal_name] = true
		}
	}

	return week_plan
}

// A big function to hold all the hand-written rules for now
// So tally things like cooking time, frequencies of dishes,
// complex things during the week, etc. and then score accordingly
func calculateScore(week_plan []database.Meal, config utilities.Config) float64 {
	// Higher numbers correspond to days where there is less time to cook
	cooking_time_score := 0.0
	duplicate_score := 0.0
	final_meal_plan_score := 0.0
	lunch_only_score := 0.0

	// Score for cooking times on days according to penalties
	for i := 0; i < len(week_plan); i++ {
		cooking_time_score += float64(week_plan[i].Cooking_time) * config.Day_weights[i]
	}

	// Penalise duplicate categories within the same week
	tmp_week_plan := make([]database.Meal, len(week_plan))
	copy(tmp_week_plan, week_plan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmp_week_plan); i++ {
		if visited[tmp_week_plan[i].Category] {
			duplicate_score += config.Duplicate_penalty
		} else {
			visited[tmp_week_plan[i].Category] = true
		}
	}

	// Penalise using lunch-only options for now
	for i := 0; i < len(week_plan); i++ {
		if week_plan[i].Lunch_only {
			lunch_only_score += config.Lunch_penalty
		}
	}

	final_meal_plan_score = cooking_time_score + duplicate_score + lunch_only_score
	return final_meal_plan_score
}

func make_meal_map(all_meals_from_database []database.Meal) map[int]database.Meal {
	meal_map := make(map[int]database.Meal)
	for i := 0; i < len(all_meals_from_database); i++ {
		meal_map[all_meals_from_database[i].ID] = all_meals_from_database[i]
	}
	return meal_map
}
