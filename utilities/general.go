package utilities

import (
	"errors"
	"fmt"
	"sort"

	"github.com/golang/glog"
	"github.com/roberto-aldera/meal-planner/database"
)

func MakeMealMap(allMealsFromDatabase []database.Meal) (mealMap map[int]database.Meal, err error) {
	if len(allMealsFromDatabase) > 0 {
		mealMap = make(map[int]database.Meal)
		for i := 0; i < len(allMealsFromDatabase); i++ {
			mealMap[allMealsFromDatabase[i].ID] = allMealsFromDatabase[i]
		}
	} else {
		err = errors.New("no meals in database")
	}
	return mealMap, err
}

// Return a slice that is partially filled by the requests.
// Also edit the meal map here, to delete requested meals as viable options.
// Note maps are passed by reference, so this function will edit the mealMap (and not return a new one).
func LoadMealRequestsAndUpdateMap(mealMap map[int]database.Meal, config Config) (weekPlanWithRequests []database.Meal,
	err error) {
	weekPlanWithRequests = make([]database.Meal, 7)

	for idx, weekDay := range config.PreferenceMealDaysOfWeek {
		weekPlanWithRequests[weekDay] = mealMap[config.PreferenceMealIDs[idx]]
		delete(mealMap, config.PreferenceMealIDs[idx])
	}
	return weekPlanWithRequests, err
}

// Note maps are passed by reference, so this function will edit the mealMap (and not return a new one).
func RemoveSpecificMeals(mealMap map[int]database.Meal, mealsToExclude []int) (err error) {
	for _, item := range mealsToExclude {
		_, keyIsValid := mealMap[item]
		if keyIsValid {
			glog.Info("Removing ", mealMap[item].MealName)
			delete(mealMap, item)
		} else {
			err = fmt.Errorf("meal key doesn't exist: %d", item)
		}
	}
	return err
}
func GetMealCategories(mealMap map[int]database.Meal) (categories []string, err error) {
	for _, meal := range mealMap {
		if !IsInSlice(categories, meal.Category) {
			categories = append(categories, meal.Category)
		}
	}
	if len(categories) < 1 {
		err = errors.New("no categories found in meal database")
	}
	// sort categories to ensure order is always the same (iterating over map is non-deterministic)
	sort.Strings(categories)
	return categories, err
}
func GetMealsInCategory(category string, mealMap map[int]database.Meal) (mealsInCategory []int, err error) {
	mealsInCategory = make([]int, 0)
	for _, meal := range mealMap {
		if meal.Category == category {
			mealsInCategory = append(mealsInCategory, meal.ID)
		}
	}
	if len(mealsInCategory) == 0 {
		// Then this category wasn't found in the database which shouldn't happen
		err = fmt.Errorf("requested category %s not found in database", category)
	}
	return mealsInCategory, err
}

func GetLunchMeals(is_lunch bool, mealMap map[int]database.Meal) (lunchMeals []int, err error) {
	lunchMeals = make([]int, 0)
	for _, meal := range mealMap {
		if meal.LunchOnly {
			lunchMeals = append(lunchMeals, meal.ID)
		}
	}
	if len(lunchMeals) == 0 {
		// Then no lunch meals were found in the database which we must flag
		err = errors.New("no lunch meals found in database")
	}
	return lunchMeals, err
}

func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
