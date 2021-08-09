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

func MakeMealPlan() {
	log.Println("Running policy...")

	// Load meals from database and print out all candidates
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()
	all_meals_from_database := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	meal_map := makeMealMap(all_meals_from_database)
	categories := utilities.GetMealCategories(meal_map)
	utilities.PrintMealDatabaseWithCategories(all_meals_from_database, categories)

	config := utilities.LoadConfiguration("meal_planner_config.json")

	if !utilities.ValidateConfiguration(config) {
		fmt.Println("Configuration is invalid!")
	}
	week_plan_with_requests, meal_map := loadMealRequestsAndUpdateMap(meal_map, config)
	meal_map = removeSpecialItems(meal_map, config.Special_exclusions, config.Previous_meals_to_exclude)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Your requested meals:")
	utilities.PrintMealPlan(week_plan_with_requests)
	fmt.Println("--------------------------------------------------------------------------------")

	best_score := config.Minimum_score // lower is better
	best_meal_plan := make([]database.Meal, len(week_plan_with_requests))
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < config.Number_of_iterations; i++ {
		week_plan := pickRandomMealsWithMap(meal_map, week_plan_with_requests, config)
		meal_plan_score := utilities.CalculateScore(week_plan, config)
		if meal_plan_score < best_score {
			// fmt.Println("New high score:", meal_plan_score, "idx = ", i)
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

func pickRandomMealsWithMap(meal_map map[int]database.Meal, week_plan_with_requests []database.Meal, config utilities.Config) []database.Meal {
	// Store map keys in a slice, and get N random items from this slice to use in the plan (to avoid picking duplicates)
	slice_of_keys := make([]int, 0)
	for key := range meal_map {
		slice_of_keys = append(slice_of_keys, key)
	}

	// Get random subset of meals to store
	random_indices := rand.Perm(len(meal_map))
	key_subset := make([]int, 0)
	for i := 0; i < len(week_plan_with_requests); i++ {
		key_subset = append(key_subset, slice_of_keys[random_indices[i]])
	}

	// Insert stored meals into week plan
	week_plan := make([]database.Meal, len(week_plan_with_requests))
	copy(week_plan, week_plan_with_requests)
	for idx := 0; idx < len(week_plan); idx++ {
		if week_plan[idx].ID == 0 { // indicates an empty slot in the week plan that can be filled
			meal_under_test := meal_map[key_subset[idx]] // get a proposed meal
			week_plan[idx] = meal_under_test
		}
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

func makeMealMap(all_meals_from_database []database.Meal) map[int]database.Meal {
	meal_map := make(map[int]database.Meal)
	for i := 0; i < len(all_meals_from_database); i++ {
		meal_map[all_meals_from_database[i].ID] = all_meals_from_database[i]
	}
	return meal_map
}

// Return a slice that is partially filled by the requests
// Possibly also edit the meal map here, to delete reuqested meals as viable options?
// Maybe that's better in another function that is called just after this one.
func loadMealRequestsAndUpdateMap(meal_map map[int]database.Meal, config utilities.Config) ([]database.Meal, map[int]database.Meal) {
	week_plan_with_requests := make([]database.Meal, 7)

	// Quick check that the inputs are legal, which really should be done in a config validation somewhere...
	if len(config.Preference_meal_IDs) == len(config.Preference_meal_days_of_week) {
		for idx, week_day := range config.Preference_meal_days_of_week {
			week_plan_with_requests[week_day] = meal_map[config.Preference_meal_IDs[idx]]
			delete(meal_map, config.Preference_meal_IDs[idx])
		}
	}
	return week_plan_with_requests, meal_map
}

// Remove meals that are to never be automatically included (like going out for dinner)
func removeSpecialItems(meal_map map[int]database.Meal, special_exclusions []int, previous_meals_to_exclude []int) map[int]database.Meal {
	for _, item := range special_exclusions {
		delete(meal_map, item)
	}
	for _, item := range previous_meals_to_exclude {
		delete(meal_map, item)
	}
	return meal_map
}
