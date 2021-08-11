package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/roberto-aldera/meal-planner/database"
)

type Config struct {
	NumberOfIterations       int
	DayWeights               [7]float64
	MinimumScore             float64
	DuplicatePenalty         float64
	LunchPenalty             float64
	PreferenceMealIDs        []int
	PreferenceMealDaysOfWeek []int
	PreviousMealsToExclude   []int
	SpecialExclusions        []int
	ExcludeSoups             bool
}

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

func LoadConfiguration(configFilePath string) Config {
	var configuration Config
	file, readErr := os.Open(configFilePath)
	if readErr != nil {
		fmt.Println(readErr.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoderErr := decoder.Decode(&configuration)
	if decoderErr != nil {
		fmt.Println("error:", decoderErr)
	}
	return configuration
}

func ValidateConfiguration(configuration Config) {
	if configuration.NumberOfIterations < 1 || configuration.NumberOfIterations > 1000000 {
		errorString := fmt.Sprintf("Configuration error. Number of iterations is is outside of range: %d", configuration.NumberOfIterations)
		panic(errorString)
	}
	for _, weight := range configuration.DayWeights {
		if weight < -100 || weight > 100 {
			errorString := fmt.Sprintf("Configuration error. Day weight is unreasonable: %f", weight)
			panic(errorString)
		}
	}
	if configuration.MinimumScore < 0 {
		errorString := fmt.Sprintf("Configuration error. Minimum score is negative: %f", configuration.MinimumScore)
		panic(errorString)
	}
	if configuration.DuplicatePenalty < 0 {
		errorString := fmt.Sprintf("Configuration error. Duplicate penalty is negative: %f", configuration.DuplicatePenalty)
		panic(errorString)
	}
	if configuration.LunchPenalty < 0 {
		errorString := fmt.Sprintf("Configuration error. Lunch penalty is negative: %f", configuration.LunchPenalty)
		panic(errorString)
	}
	if len(configuration.PreferenceMealIDs) < 0 || len(configuration.PreferenceMealIDs) > 7 {
		panic("Configuration error. PreferenceMealIDs length is out of range (0,7)")
	}
	if len(configuration.PreferenceMealDaysOfWeek) < 0 || len(configuration.PreferenceMealDaysOfWeek) > 7 {
		panic("Configuration error. PreferenceMealDaysOfWeek length is out of range")
	}
	if len(configuration.PreferenceMealIDs) != len(configuration.PreferenceMealDaysOfWeek) {
		panic("Configuration error. PreferenceMealIDs length is different to PreferenceMealDaysOfWeek")
	}
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

func IsInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
