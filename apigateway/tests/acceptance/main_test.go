package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/ProyectoFinal-Microservicios/architecture/apigateway/tests/acceptance/steps"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
	Paths:  []string{"features"},
	Tags:   "",
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	status := godog.TestSuite{
		Name: "API Gateway Acceptance Tests",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			apiCtx := steps.NewAPIGatewayContext()
			steps.InitializeScenario(ctx, apiCtx)
		},
		Options: &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("test suite failed with status: %d", status)
	}
}

func main() {
	flag.Parse()

	status := godog.TestSuite{
		Name: "API Gateway Acceptance Tests",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			apiCtx := steps.NewAPIGatewayContext()
			steps.InitializeScenario(ctx, apiCtx)
		},
		Options: &opts,
	}.Run()

	if st := os.Getenv("GODOG_PUBLISH_QUIET"); st != "off" {
		fmt.Println("Test execution completed with status:", status)
	}

	os.Exit(status)
}