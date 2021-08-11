package strategy

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/glog"

	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/utilities"
)

func MakeMealPlan() {

	configFilePath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	fmt.Println("Running policy...")

	// Load meals from database and print out all candidates
	sqliteDatabase, _ := sql.Open("sqlite3", "/Users/roberto/github-code/meal-planner/localdata/meal-data.db")
	defer sqliteDatabase.Close()
	allMealsFromDatabase := database.LoadDatabaseEntriesIntoContainer(sqliteDatabase)

	mealMap := makeMealMap(allMealsFromDatabase)
	categories := utilities.GetMealCategories(mealMap)
	utilities.PrintMealDatabaseWithCategories(allMealsFromDatabase, categories)

	config := utilities.LoadConfiguration(*configFilePath)

	utilities.ValidateConfiguration(config)

	weekPlanWithRequests, mealMap := loadMealRequestsAndUpdateMap(mealMap, config)
	utilities.PrintExcludedMeals(mealMap, config.PreviousMealsToExclude)
	mealMap = removeSpecificMeals(mealMap, config.SpecialExclusions)
	mealMap = removeSpecificMeals(mealMap, config.PreviousMealsToExclude)
	if config.ExcludeSoups {
		soups := getMealsInCategory("Soups", mealMap)
		mealMap = removeSpecificMeals(mealMap, soups)
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

func makeMealMap(allMealsFromDatabase []database.Meal) map[int]database.Meal {
	mealMap := make(map[int]database.Meal)
	for i := 0; i < len(allMealsFromDatabase); i++ {
		mealMap[allMealsFromDatabase[i].ID] = allMealsFromDatabase[i]
	}
	return mealMap
}

// Return a slice that is partially filled by the requests
// Possibly also edit the meal map here, to delete reuqested meals as viable options?
// Maybe that's better in another function that is called just after this one.
func loadMealRequestsAndUpdateMap(mealMap map[int]database.Meal, config utilities.Config) ([]database.Meal, map[int]database.Meal) {
	weekPlanWithRequests := make([]database.Meal, 7)

	// Quick check that the inputs are legal, which really should be done in a config validation somewhere...
	if len(config.PreferenceMealIDs) == len(config.PreferenceMealDaysOfWeek) {
		for idx, weekDay := range config.PreferenceMealDaysOfWeek {
			weekPlanWithRequests[weekDay] = mealMap[config.PreferenceMealIDs[idx]]
			delete(mealMap, config.PreferenceMealIDs[idx])
		}
	}
	return weekPlanWithRequests, mealMap
}

func removeSpecificMeals(mealMap map[int]database.Meal, mealsToExclude []int) map[int]database.Meal {
	for _, item := range mealsToExclude {
		_, keyIsValid := mealMap[item]
		if keyIsValid {
			glog.Info("Removing ", mealMap[item].MealName)
			delete(mealMap, item)
		} else {
			panic(fmt.Sprintf("Meal key doesn't exist: %d", item))
		}
	}
	return mealMap
}

func getMealsInCategory(category string, mealMap map[int]database.Meal) []int {
	// TODO: validate that category is correct (it must exist)
	mealsInCategory := make([]int, 0)
	for _, meal := range mealMap {
		if meal.Category == category {
			mealsInCategory = append(mealsInCategory, meal.ID)
		}
	}
	return mealsInCategory

}
