package utilities

import (
	"errors"
	"fmt"

	"github.com/roberto-aldera/meal-planner/database"
)

func PrintMealDatabase(mealDatabase []database.Meal) (err error) {
	if len(mealDatabase) == 0 {
		err = errors.New("meal database is empty")
		return err
	}
	fmt.Println("Meals available are:")
	for _, meal := range mealDatabase {
		fmt.Println(meal.ID, "->", meal.MealName)
	}
	return err
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
	dayNames := [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if len(weekPlan) == 7 {
		for i := range dayNames {
			if weekPlan[i].ID != 0 {
				fmt.Printf("%s: \t %s -> %d \n", dayNames[i], weekPlan[i].MealName, weekPlan[i].ID)
			} else {
				fmt.Printf("%s: \n", dayNames[i])
			}
		}
	} else {
		panic("Meal plan not complete. Expected to be of length 7.")
	}
}
