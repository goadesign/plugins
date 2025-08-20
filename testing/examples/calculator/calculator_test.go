package calculatorapi

import (
	"testing"

	"goa.design/plugins/v3/testing/examples/calculator/gen/calculator/calculatortest"
)

// Use the generated test suite
func TestCalculator(t *testing.T) {
	svc := NewCalculator()
	RunCalculatorHarness(t, svc)
}

func TestCalculatorScenarios(t *testing.T) {
	// Create test harness
	service := NewCalculator()
	h := calculatortest.NewHarness(t, service)
	defer h.Close()
	
	// Load scenarios from YAML
	runner, err := calculatortest.LoadScenarios("scenarios.yaml")
	if err != nil {
		t.Fatalf("failed to load scenarios: %v", err)
	}
	
	// Run all scenarios
	runner.Run(t, h.Client)
}

// TestSpecificScenario demonstrates running a single scenario
func TestSpecificScenario(t *testing.T) {
	// Create test harness
	service := NewCalculator()
	h := calculatortest.NewHarness(t, service)
	defer h.Close()
	
	// Load scenarios
	runner, err := calculatortest.LoadScenarios("scenarios.yaml")
	if err != nil {
		t.Fatalf("failed to load scenarios: %v", err)
	}
	
	// Run specific scenario
	runner.RunNamed(t, h.Client, "division_with_validation")
}