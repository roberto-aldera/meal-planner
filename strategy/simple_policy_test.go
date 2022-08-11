package strategy

import (
	"testing"
)

func TestMakeMealPlan(t *testing.T) {
	configFilePath := "../default_config.json"
	MakeMealPlan(configFilePath)
}
