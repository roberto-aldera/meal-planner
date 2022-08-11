package utilities

import (
	"testing"
)

// Check that configuration loads as expected
func TestLoadConfiguration(t *testing.T) {
	filePath := "nonsense.json"
	_, readErr := LoadConfiguration(filePath)
	if readErr == nil {
		t.Fatal("Expected error to not be nil as configuration file path doesn't exist.")
	}

	filePath = "../default_config.json"
	_, readErr = LoadConfiguration(filePath)
	if readErr != nil {
		t.Fatal(readErr)
	}
}

func TestValidateConfiguration(t *testing.T) {
	configuration, _ := LoadConfiguration("../default_config.json")

	// Check first that all runs as expected
	err := ValidateConfiguration(configuration)
	if err != nil {
		t.Fatal(err)
	}

	bad_configuration := configuration
	configuration_error_string := "Expected bad configuration to be rejected."

	run_config_validation := func(bad_configuration Config) {
		err = ValidateConfiguration(bad_configuration)
		if err == nil {
			t.Fatal(configuration_error_string)
		}
	}

	bad_configuration.NumberOfIterations = 0
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.NumberOfIterations = 1e9
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.ComplexMealRequested = bad_configuration.ComplexMealRequested[1:]
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.SimpleMealRequested = bad_configuration.SimpleMealRequested[1:]
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	requests := []bool{true, false, false, false, false, false, false}
	bad_configuration.SimpleMealRequested = requests
	bad_configuration.ComplexMealRequested = requests
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.MinimumScore = -1
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.ScorePenalty = -1
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.PreferenceMealIDs = []int{200, 201, 202, 203, 204, 205, 206, 207}
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.PreferenceMealDaysOfWeek = []int{0, 1, 2, 3, 4, 5, 6, 7}
	run_config_validation(bad_configuration)

	bad_configuration = configuration
	bad_configuration.PreferenceMealIDs = []int{200}
	bad_configuration.PreferenceMealDaysOfWeek = []int{0, 1}
	run_config_validation(bad_configuration)

}
