package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	MealDatabasePath                  string
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

func LoadConfiguration(configFilePath string) (configuration Config, err error) {
	file, readErr := os.Open(configFilePath)
	if readErr != nil {
		fmt.Println(readErr.Error())
	}
	defer file.Close()
	return decodeConfiguration(file)
}

func decodeConfiguration(handle io.Reader) (configuration Config, err error) {
	decoder := json.NewDecoder(handle)
	decoderErr := decoder.Decode(&configuration)
	return configuration, decoderErr
}

func ValidateConfiguration(configuration Config) (err error) {
	if configuration.NumberOfIterations < 1 || configuration.NumberOfIterations > 1000000 {
		return fmt.Errorf("number of iterations is is outside of range: %d", configuration.NumberOfIterations)
	}
	fmt.Println(len(configuration.ComplexMealRequested), configuration.ComplexMealRequested)
	if len(configuration.ComplexMealRequested) != 7 {
		return fmt.Errorf("complexMealRequested length must be 7, got %d", len(configuration.ComplexMealRequested))
	}
	if len(configuration.SimpleMealRequested) != 7 {
		return fmt.Errorf("simpleMealRequested length must be 7")
	}
	for idx := range configuration.ComplexMealRequested {
		if configuration.ComplexMealRequested[idx] && configuration.SimpleMealRequested[idx] {
			return fmt.Errorf("cannot request both complex and simple meal on the same day for index = %d", idx)
		}
	}
	if configuration.MinimumScore < 0 {
		return fmt.Errorf("minimum score is negative: %f", configuration.MinimumScore)
	}
	if configuration.ScorePenalty < 0 {
		return fmt.Errorf("duplicate penalty is negative: %f", configuration.ScorePenalty)
	}
	if len(configuration.PreferenceMealIDs) < 0 || len(configuration.PreferenceMealIDs) > 7 {
		return fmt.Errorf("preferenceMealIDs length is out of range (0,7)")
	}
	if len(configuration.PreferenceMealDaysOfWeek) < 0 || len(configuration.PreferenceMealDaysOfWeek) > 7 {
		return fmt.Errorf("preferenceMealDaysOfWeek length is out of range")
	}
	if len(configuration.PreferenceMealIDs) != len(configuration.PreferenceMealDaysOfWeek) {
		return fmt.Errorf("preferenceMealIDs length is different to PreferenceMealDaysOfWeek")
	}
	return nil
}
