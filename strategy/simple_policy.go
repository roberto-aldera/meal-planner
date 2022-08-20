package strategy

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/roberto-aldera/meal-planner/database"
	"github.com/roberto-aldera/meal-planner/utilities"
)

func MakeMealPlan(config utilities.Config, allMealsFromDatabase []database.Meal) (err error) {
	mealMap, err := setupMealMap(config, allMealsFromDatabase)
	if err != nil {
		return fmt.Errorf("setupMealMap has failed %s", err)
	}

	weekPlanWithRequests, err := assignPreferences(config, mealMap)
	if err != nil {
		return fmt.Errorf("assignPreferences has failed %s", err)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Your requested meals:")
	err = utilities.PrintMealPlan(weekPlanWithRequests)
	if err != nil {
		fmt.Printf("PrintMealPlan failed: %s", err)
	}
	fmt.Println("--------------------------------------------------------------------------------")

	bestScore := config.MinimumScore // lower is better
	bestMealPlan := make([]database.Meal, len(weekPlanWithRequests))
	// rand.Seed(1624728791619452000) // hardcoded for easier debugging
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < config.NumberOfIterations; i++ {
		weekPlan := pickRandomMealsWithMap(mealMap, weekPlanWithRequests, config)
		mealPlanScore, err := utilities.CalculateScore(weekPlan, config)
		if err != nil {
			fmt.Printf("CalculateScore failed: %s", err)
		}
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
	return err
}

func setupMealMap(config utilities.Config, allMealsFromDatabase []database.Meal) (map[int]database.Meal, error) {
	mealMap, err := utilities.MakeMealMap(allMealsFromDatabase)
	if err != nil {
		fmt.Printf("MakeMealMap has failed: %s", err.Error())
		return nil, err
	}
	categories, err := utilities.GetMealCategories(mealMap)
	if err != nil {
		fmt.Printf("GetMealCategories has failed: %s", err.Error())
		return nil, err
	}

	fmt.Println("Meals available are:")
	for _, category := range categories {
		fmt.Println("\n------------------------------>", category)
		for _, meal := range allMealsFromDatabase {
			if meal.Category == category {
				fmt.Println(meal.ID, "->", meal.MealName)
			}
		}
	}
	fmt.Println("\n--------------------------------------------------------------------------------")

	return mealMap, err
}

func assignPreferences(config utilities.Config, mealMap map[int]database.Meal) ([]database.Meal, error) {
	weekPlanWithRequests, err := utilities.LoadMealRequestsAndUpdateMap(mealMap, config)
	if err != nil {
		fmt.Printf("LoadMealRequestsAndUpdateMap failed: %s", err)
		return nil, err
	}

	// Print any meal exclusions
	if (len(config.PreviousMealsToExclude)) > 0 {
		fmt.Println("These meals have been requested to be excluded:")
		for _, mealID := range config.PreviousMealsToExclude {
			_, keyIsValid := mealMap[mealID]
			if keyIsValid {
				fmt.Println(mealMap[mealID].MealName, "->", mealMap[mealID].ID)
			} else {
				return nil, fmt.Errorf("meal ID doesn't exist: %d", mealID)
			}
		}
	} else {
		fmt.Println("No meals were requested to be excluded.")
	}

	err = utilities.RemoveSpecificMeals(mealMap, config.SpecialExclusions)
	if err != nil {
		return nil, fmt.Errorf("RemoveSpecificMeals failed: %s", err)
	}
	if config.ExcludeSoups {
		soups, err := utilities.GetMealsInCategory("Soups", mealMap)
		if err != nil {
			return nil, fmt.Errorf("GetMealsInCategory failed: %s", err)
		}
		// no need to check error, because we literally just got soups from the mealMap
		utilities.RemoveSpecificMeals(mealMap, soups)
	}
	if config.ExcludeLunches {
		lunches, err := utilities.GetLunchMeals(mealMap)
		if err != nil {
			return nil, fmt.Errorf("GetLunchMeals error: %s", err)
		}
		// no need to check error, because we literally just got lunches from the mealMap
		utilities.RemoveSpecificMeals(mealMap, lunches)
	}
	return weekPlanWithRequests, err
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
