package utilities

import (
	"fmt"

	"github.com/roberto-aldera/meal-planner/database"
)

type Config struct {
	Number_of_iterations         int
	Day_weights                  [7]float64
	Minimum_score                float64
	Duplicate_penalty            float64
	Lunch_penalty                float64
	Preference_meal_IDs          []int
	Preference_meal_days_of_week []int
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

func PrintMealDatabaseWithCategories(meal_database []database.Meal, categories []string) {
	fmt.Println("Meals available are:")
	for _, category := range categories {
		fmt.Println("\n------------------------------>", category)

		for _, meal := range meal_database {
			if meal.Category == category {
				fmt.Println(meal.ID, "->", meal.Meal_name)
			}
		}
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

func ValidateConfiguration(configuration Config) bool {
	// TODO: use panic perhaps?
	fmt.Println("Running configuration validation...")

	if configuration.Number_of_iterations < 1 || configuration.Number_of_iterations > 1000000 {
		fmt.Println("Configuration error: number of iterations", configuration.Number_of_iterations, "is outside of range")
		return false
	}
	for _, weight := range configuration.Day_weights {
		if weight < -100 || weight > 100 {
			fmt.Println("Configuration error: day weight of", weight, "is unreasonable")
			return false
		}
	}
	if configuration.Minimum_score < 0 {
		fmt.Println("Configuration error: minimum score", configuration.Minimum_score, "is negative")
		return false
	}
	if configuration.Duplicate_penalty < 0 {
		fmt.Println("Configuration error: duplicate penalty", configuration.Duplicate_penalty, "is negative")
		return false
	}
	if configuration.Lunch_penalty < 0 {
		fmt.Println("Configuration error: lunch penalty", configuration.Lunch_penalty, "is negative")
		return false
	}
	if len(configuration.Preference_meal_IDs) < 0 || len(configuration.Preference_meal_IDs) > 7 {
		fmt.Println("Configuration error: Preference_meal_IDs length is out of range")
		return false
	}
	if len(configuration.Preference_meal_days_of_week) < 0 || len(configuration.Preference_meal_days_of_week) > 7 {
		fmt.Println("Configuration error: Preference_meal_days_of_week length is out of range")
		return false
	}
	if len(configuration.Preference_meal_IDs) != len(configuration.Preference_meal_days_of_week) {
		fmt.Println("Configuration error: Preference_meal_IDs length is different to Preference_meal_days_of_week")
		return false
	}

	return true
}

func GetMealCategories(meal_map map[int]database.Meal) []string {
	categories := make([]string, 0)
	for _, meal := range meal_map {
		if !IsInSlice(categories, meal.Category) {
			categories = append(categories, meal.Category)
		}
	}
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
