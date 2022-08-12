package strategy

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/utilities"
)

func MakeMealPlan(configFilePath string) {

	fmt.Println("Running policy...")

	// Load meals from database and print out all candidates
	sqliteDatabase, _ := sql.Open("sqlite3", "../meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	mealMap, err := utilities.MakeMealMap(allMealsFromDatabase)
	if err != nil {
		fmt.Println(err.Error())
	}

	categories, err := utilities.GetMealCategories(mealMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	utilities.PrintMealDatabaseWithCategories(allMealsFromDatabase, categories)

	config, err := utilities.LoadConfiguration(configFilePath)
	if err != nil {
		panic("Configuration has failed to load.")
	}

	err = utilities.ValidateConfiguration(config)
	if err != nil {
		panic(fmt.Sprintf("Configuration validation failed: %s", err))
	}

	weekPlanWithRequests, mealMap := utilities.LoadMealRequestsAndUpdateMap(mealMap, config)
	utilities.PrintExcludedMeals(mealMap, config.PreviousMealsToExclude)
	mealMap = utilities.RemoveSpecificMeals(mealMap, config.SpecialExclusions)
	mealMap = utilities.RemoveSpecificMeals(mealMap, config.PreviousMealsToExclude)
	if config.ExcludeSoups {
		soups := utilities.GetMealsInCategory("Soups", mealMap)
		mealMap = utilities.RemoveSpecificMeals(mealMap, soups)
	}
	if config.ExcludeLunches {
		lunches := utilities.GetLunchMeals(config.ExcludeLunches, mealMap)
		mealMap = utilities.RemoveSpecificMeals(mealMap, lunches)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Your requested meals:")
	utilities.PrintMealPlan(weekPlanWithRequests)
	fmt.Println("--------------------------------------------------------------------------------")

	bestScore := config.MinimumScore // lower is better
	bestMealPlan := make([]database.Meal, len(weekPlanWithRequests))
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < config.NumberOfIterations; i++ {
		weekPlan := pickRandomMealsWithMap(mealMap, weekPlanWithRequests, config)
		mealPlanScore := utilities.CalculateScore(weekPlan, config)
		if mealPlanScore < bestScore {
			bestMealPlan = weekPlan
			bestScore = mealPlanScore
		}
	}

	if len(bestMealPlan) == 7 {
		fmt.Println("Best meal plan after", config.NumberOfIterations, "iterations from a total of", len(allMealsFromDatabase), "meals:")
		utilities.PrintMealPlan(bestMealPlan)
		fmt.Println("Score:", bestScore)
	} else {
		fmt.Println("No valid meal plan was possible with the provided requirements.")
	}
}

func pickRandomMealsWithMap(mealMap map[int]database.Meal, weekPlanWithRequests []database.Meal, config utilities.Config) []database.Meal {
	// Store map keys in a slice, and get N random items from this slice to use in the plan (to avoid picking duplicates)
	sliceOfKeys := make([]int, 0)
	for key := range mealMap {
		sliceOfKeys = append(sliceOfKeys, key)
	}

	// Get random subset of meals to store
	randomIndices := rand.Perm(len(mealMap))
	keySubset := make([]int, 0)
	for i := 0; i < len(weekPlanWithRequests); i++ {
		keySubset = append(keySubset, sliceOfKeys[randomIndices[i]])
	}

	// Insert stored meals into week plan
	weekPlan := make([]database.Meal, len(weekPlanWithRequests))
	copy(weekPlan, weekPlanWithRequests)
	for idx := 0; idx < len(weekPlan); idx++ {
		if weekPlan[idx].ID == 0 { // indicates an empty slot in the week plan that can be filled
			mealUnderTest := mealMap[keySubset[idx]] // get a proposed meal
			weekPlan[idx] = mealUnderTest
		}
	}

	// Debug: check for duplicates
	tmpWeekPlan := make([]database.Meal, len(weekPlan))
	copy(tmpWeekPlan, weekPlan)
	visited := make(map[string]bool)
	for i := 0; i < len(tmpWeekPlan); i++ {
		if visited[tmpWeekPlan[i].MealName] {
			fmt.Println("*** Duplicate found:", tmpWeekPlan[i].MealName)
		} else {
			visited[tmpWeekPlan[i].MealName] = true
		}
	}

	return weekPlan
}
