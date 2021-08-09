package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

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
	Previous_meals_to_exclude    []int
	Special_exclusions           []int
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

func LoadConfiguration(config_file_path string) Config {
	var configuration Config
	file, read_err := os.Open(config_file_path)
	if read_err != nil {
		fmt.Println(read_err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder_err := decoder.Decode(&configuration)
	if decoder_err != nil {
		fmt.Println("error:", decoder_err)
	}
	return configuration
}

func ValidateConfiguration(configuration Config) {
	if configuration.Number_of_iterations < 1 || configuration.Number_of_iterations > 1000000 {
		errorString := fmt.Sprintf("Configuration error. Number of iterations is is outside of range: %d", configuration.Number_of_iterations)
		panic(errorString)
	}
	for _, weight := range configuration.Day_weights {
		if weight < -100 || weight > 100 {
			errorString := fmt.Sprintf("Configuration error. Day weight is unreasonable: %f", weight)
			panic(errorString)
		}
	}
	if configuration.Minimum_score < 0 {
		errorString := fmt.Sprintf("Configuration error. Minimum score is negative: %f", configuration.Minimum_score)
		panic(errorString)
	}
	if configuration.Duplicate_penalty < 0 {
		errorString := fmt.Sprintf("Configuration error. Duplicate penalty is negative: %f", configuration.Duplicate_penalty)
		panic(errorString)
	}
	if configuration.Lunch_penalty < 0 {
		errorString := fmt.Sprintf("Configuration error. Lunch penalty is negative: %f", configuration.Lunch_penalty)
		panic(errorString)
	}
	if len(configuration.Preference_meal_IDs) < 0 || len(configuration.Preference_meal_IDs) > 7 {
		panic("Configuration error. Preference_meal_IDs length is out of range (0,7)")
	}
	if len(configuration.Preference_meal_days_of_week) < 0 || len(configuration.Preference_meal_days_of_week) > 7 {
		panic("Configuration error. Preference_meal_days_of_week length is out of range")
	}
	if len(configuration.Preference_meal_IDs) != len(configuration.Preference_meal_days_of_week) {
		panic("Configuration error. Preference_meal_IDs length is different to Preference_meal_days_of_week")
	}
}

func GetMealCategories(meal_map map[int]database.Meal) []string {
	categories := make([]string, 0)
	for _, meal := range meal_map {
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
