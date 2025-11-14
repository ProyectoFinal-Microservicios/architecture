package acceptance

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/ProyectoFinal-Microservicios/architecture/apigateway/tests/acceptance/steps"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "API Gateway Acceptance Tests",
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}

	if status := suite.Run(); status != 0 {
		t.Fatalf("test suite failed with status: %d", status)
	}
}