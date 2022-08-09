package utilities

import (
	"testing"
)

// Check that configuration loads as expected
func TestLoadConfiguration(t *testing.T) {
	filePath := "nonsense.json"
	_, readErr := LoadConfiguration(filePath)
	if readErr == nil {
		t.Fatal(readErr)
	}

	filePath = "../default_config.json"
	_, readErr = LoadConfiguration(filePath)
	if readErr != nil {
		t.Fatal(readErr)
	}
}
