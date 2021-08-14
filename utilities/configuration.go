package utilities

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NumberOfIterations                int
	ComplexMealRequested              []bool
	SimpleMealRequested               []bool
	MinimumScore                      float64
	ScorePenalty                      float64
	DefinitionOfLongMealPrepTimeHours float64
	PreferenceMealIDs                 []int
	PreferenceMealDaysOfWeek          []int
	PreviousMealsToExclude            []int
	SpecialExclusions                 []int
	ExcludeLunches                    bool
	ExcludeSoups                      bool
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
	fmt.Println(len(configuration.ComplexMealRequested), configuration.ComplexMealRequested)
	if len(configuration.ComplexMealRequested) != 7 {
		panic("Configuration error. ComplexMealRequested length must be 7.")
	}
	if len(configuration.SimpleMealRequested) != 7 {
		panic("Configuration error. SimpleMealRequested length must be 7.")
	}
	for idx := range configuration.ComplexMealRequested {
		if configuration.ComplexMealRequested[idx] && configuration.SimpleMealRequested[idx] {
			errorString := fmt.Sprintf("Configuration error. Cannot request both complex and simple meal on the same day for index = %d", idx)
			panic(errorString)
		}
	}
	if configuration.MinimumScore < 0 {
		errorString := fmt.Sprintf("Configuration error. Minimum score is negative: %f", configuration.MinimumScore)
		panic(errorString)
	}
	if configuration.ScorePenalty < 0 {
		errorString := fmt.Sprintf("Configuration error. Duplicate penalty is negative: %f", configuration.ScorePenalty)
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
