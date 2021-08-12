package utilities

import (
	"sort"

	"github.com/roberto-aldera/meal-planner/database"
)

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

func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
