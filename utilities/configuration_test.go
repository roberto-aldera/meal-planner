package utilities

import (
	"testing"
)

// Check that configuration loads as expected
func TestLoadConfiguration(t *testing.T) {
	filePath := "../default_config.json"
	LoadConfiguration(filePath)
}
