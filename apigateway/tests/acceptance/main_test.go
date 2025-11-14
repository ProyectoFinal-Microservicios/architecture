package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/ProyectoFinal-Microservicios/architecture/apigateway/tests/acceptance/steps"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
	Paths:  []string{"tests/acceptance/features"},
	Tags:   "",
}

func init() {
	godog.BindCommandLineFlags(&opts)
}

func TestFeatures(t interface{ Errorf(string, ...interface{}) }) {
	suite := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options:             &opts,
	}

	if status := suite.Run(); status != 0 {
		t.Errorf("test suite failed with status: %d", status)
	}
}

func main() {
	flag.Parse()

	status := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := os.Getenv("GODOG_PUBLISH_QUIET"); st != "off" {
		fmt.Println("Test execution completed with status:", status)
	}

	os.Exit(status)
}
