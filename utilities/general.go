package utilities

import (
	"fmt"

	"github.com/roberto-aldera/meal-planner/database"
)

type Config struct {
	Number_of_iterations int
	Day_weights          [7]float64
	Minimum_score        float64
	Duplicate_penalty    float64
	Lunch_penalty        float64
}

type Specific_meal struct {
	Meal_ID_idx int
	Day_of_week int
}

func PrintMealDatabase(meal_database []database.Meal) {
	fmt.Println("Meals available are:")
	for _, meal := range meal_database {
		fmt.Println(meal.ID, "->", meal.Meal_name)
	}
}

func PrintMealPlan(week_plan []database.Meal) {
	if len(week_plan) == 7 {
		fmt.Println("Monday:   ", week_plan[0].Meal_name)
		fmt.Println("Tuesday:  ", week_plan[1].Meal_name)
		fmt.Println("Wednesday:", week_plan[2].Meal_name)
		fmt.Println("Thursday: ", week_plan[3].Meal_name)
		fmt.Println("Friday:   ", week_plan[4].Meal_name)
		fmt.Println("Saturday: ", week_plan[5].Meal_name)
		fmt.Println("Sunday:   ", week_plan[6].Meal_name)
	} else {
		fmt.Println("Meal plan not complete.")
	}
}
