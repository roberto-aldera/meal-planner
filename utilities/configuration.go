package utilities

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NumberOfIterations       int
	DayWeights               [7]float64
	MinimumScore             float64
	DuplicatePenalty         float64
	PreferenceMealIDs        []int
	PreferenceMealDaysOfWeek []int
	PreviousMealsToExclude   []int
	SpecialExclusions        []int
	ExcludeLunches           bool
	ExcludeSoups             bool
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
