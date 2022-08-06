package utilities

import (
	"fmt"
	"sort"

	"github.com/golang/glog"
	"github.com/roberto-aldera/meal-planner/database"
)

func MakeMealMap(allMealsFromDatabase []database.Meal) map[int]database.Meal {
	mealMap := make(map[int]database.Meal)
	for i := 0; i < len(allMealsFromDatabase); i++ {
		mealMap[allMealsFromDatabase[i].ID] = allMealsFromDatabase[i]
	}
	return mealMap
}

// Return a slice that is partially filled by the requests
// Possibly also edit the meal map here, to delete requested meals as viable options?
// Maybe that's better in another function that is called just after this one.
func LoadMealRequestsAndUpdateMap(mealMap map[int]database.Meal, config Config) ([]database.Meal, map[int]database.Meal) {
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

func RemoveSpecificMeals(mealMap map[int]database.Meal, mealsToExclude []int) map[int]database.Meal {
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
func GetMealCategories(mealMap map[int]database.Meal) []string {
	categories := make([]string, 0)
	for _, meal := range mealMap {
		if !IsInSlice(categories, meal.Category) {
			categories = append(categories, meal.Category)
		}
	}
	// sort categories to ensure order is always the same (iterating over map is non-deterministic)
	sort.Strings(categories)
	return categories
}
func GetMealsInCategory(category string, mealMap map[int]database.Meal) []int {
	// TODO: validate that category is correct (it must exist)
	mealsInCategory := make([]int, 0)
	for _, meal := range mealMap {
		if meal.Category == category {
			mealsInCategory = append(mealsInCategory, meal.ID)
		}
	}
	return mealsInCategory
}

func GetLunchMeals(is_lunch bool, mealMap map[int]database.Meal) []int {
	lunchMeals := make([]int, 0)
	for _, meal := range mealMap {
		if meal.LunchOnly {
			lunchMeals = append(lunchMeals, meal.ID)
		}
	}
	return lunchMeals
}

func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
