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

func PrintMealPlan(weekPlan []database.Meal) (err error) {
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
		err = fmt.Errorf("meal plan not complete. Expected to be of length 7, got %d", len(weekPlan))
	}
	return err
}
