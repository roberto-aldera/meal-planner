package main

import (
	"fmt"
	"math/rand"
)

type Meal struct {
	name_of_dish string
	tomato_base  bool
}

func newMeal(name string) *Meal {
	meal := Meal{}
	meal.name_of_dish = name
	meal.tomato_base = false
	return &meal
}

func makePastaPlan(all_meals []Meal) []Meal {
	// Order meals based on if they are tomato-based or not
	var pasta_plan []Meal // create empty pasta plan
	initial_meal_idx := 0
	pasta_plan = append(pasta_plan, all_meals[0])                                       // initialise with a first dish
	all_meals = append(all_meals[:initial_meal_idx], all_meals[initial_meal_idx+1:]...) // erase dish from the possible options
	fmt.Println("-- Initial pasta plan --")
	printMealPlan(pasta_plan)

	for len(pasta_plan) < 5 {
		idx := rand.Intn(len(all_meals))
		meal_under_test := all_meals[idx] // get a proposed meal
		// Now check if proposed meal is suitable based on if its tomato-based state is the opposite of the last selected meal
		if meal_under_test.tomato_base != pasta_plan[len(pasta_plan)-1].tomato_base {
			pasta_plan = append(pasta_plan, meal_under_test)
			all_meals = append(all_meals[:idx], all_meals[idx+1:]...) // erase meal from available options
		}
	}
	return pasta_plan
}

func printMealPlan(meal_plan []Meal) {
	for i := 0; i < len(meal_plan); i++ {
		fmt.Printf("%+v\n", meal_plan[i].name_of_dish)
	}
}

func generatePastasAndMakePlan() {
	fmt.Println("Meal planner is running...")

	// Make some meals to choose from
	meal_1 := newMeal("Pasta with peas")
	meal_1.tomato_base = false

	meal_2 := newMeal("Ragu")
	meal_2.tomato_base = true

	meal_3 := newMeal("Pasta with lentils")
	meal_3.tomato_base = false

	meal_4 := newMeal("Pasta e fagioli")
	meal_4.tomato_base = false

	meal_5 := newMeal("Amatriciana")
	meal_5.tomato_base = true

	// Load all the meals into a slice
	all_meals := make([]Meal, 5)
	all_meals[0] = *meal_1
	all_meals[1] = *meal_2
	all_meals[2] = *meal_3
	all_meals[3] = *meal_4
	all_meals[4] = *meal_5

	pasta_plan := makePastaPlan(all_meals)
	fmt.Println("-- Proposed pasta plan --")
	printMealPlan(pasta_plan)
}

func main() {
	generatePastasAndMakePlan()
}
