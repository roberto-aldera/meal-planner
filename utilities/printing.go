package utilities

import (
	"fmt"

	"github.com/roberto-aldera/meal-planner/database"
)

func PrintMealDatabase(mealDatabase []database.Meal) {
	fmt.Println("Meals available are:")
	for _, meal := range mealDatabase {
		fmt.Println(meal.ID, "->", meal.MealName)
	}
}

func PrintMealDatabaseWithCategories(mealDatabase []database.Meal, categories []string) {
	fmt.Println("Meals available are:")
	for _, category := range categories {
		fmt.Println("\n------------------------------>", category)
		for _, meal := range mealDatabase {
			if meal.Category == category {
				fmt.Println(meal.ID, "->", meal.MealName)
			}
		}
	}
	fmt.Println("\n--------------------------------------------------------------------------------")
}

func PrintExcludedMeals(mealMap map[int]database.Meal, previousMealsToExclude []int) {
	if (len(previousMealsToExclude)) > 0 {
		fmt.Println("These meals have been requested to be excluded:")
		for _, mealID := range previousMealsToExclude {
			fmt.Println(mealMap[mealID].MealName, "->", mealMap[mealID].ID)
		}
	} else {
		fmt.Println("No meals were requested to be excluded.")
	}
}

func PrintMealPlan(weekPlan []database.Meal) {
	if len(weekPlan) == 7 {
		fmt.Println("Monday:   ", weekPlan[0].MealName)
		fmt.Println("Tuesday:  ", weekPlan[1].MealName)
		fmt.Println("Wednesday:", weekPlan[2].MealName)
		fmt.Println("Thursday: ", weekPlan[3].MealName)
		fmt.Println("Friday:   ", weekPlan[4].MealName)
		fmt.Println("Saturday: ", weekPlan[5].MealName)
		fmt.Println("Sunday:   ", weekPlan[6].MealName)
	} else {
		fmt.Println("Meal plan not complete.")
	}
}
