package environment_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bztftest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnvironment_BasicById(t *testing.T) {
	var environment environments.Environment
	// Create random env name
	name := bztftest.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { bztftest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: bztftest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: environmentByIdConfig(name),
				Check: resource.ComposeTestCheckFunc(
					bztftest.TestAccCheckEnvironmentExists("data.bastionzero_environment.env", &environment),
					resource.ComposeTestCheckFunc(
						checkEnvironmentSchemaAttrs("data.bastionzero_environment.env", expectedEnvironmentSchema{ExpectedName: name})...,
					),
				),
			},
		},
	})
}

type expectedEnvironmentSchema struct {
	ExpectedName                       string
	ExpectedDescription                string
	ExpectedOfflineCleanupTimeoutHours string
	ExpectedTargetsLength              string
}

func checkEnvironmentSchemaAttrs(stateName string, expected expectedEnvironmentSchema) []resource.TestCheckFunc {
	// Default assertions
	if expected.ExpectedOfflineCleanupTimeoutHours == "" {
		expected.ExpectedOfflineCleanupTimeoutHours = strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)
	}
	if expected.ExpectedTargetsLength == "" {
		expected.ExpectedTargetsLength = "0"
	}

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(stateName, "name", expected.ExpectedName),
		resource.TestCheckResourceAttr(stateName, "description", expected.ExpectedDescription),
		resource.TestCheckResourceAttr(stateName, "offline_cleanup_timeout_hours", expected.ExpectedOfflineCleanupTimeoutHours),
		resource.TestCheckResourceAttr(stateName, "is_default", "false"),
		resource.TestCheckResourceAttr(stateName, "targets.%", expected.ExpectedTargetsLength),
		resource.TestMatchResourceAttr(stateName, "id", bztftest.ExpectedIDRegEx()),
		resource.TestMatchResourceAttr(stateName, "organization_id", bztftest.ExpectedIDRegEx()),
		resource.TestMatchResourceAttr(stateName, "time_created", bztftest.ExpectedTimestampRegEx()),
	}
}

func environmentByIdConfig(envName string) string {
	return fmt.Sprintf(`
%s
data "bastionzero_environment" "env" {
  id = bastionzero_environment.env.id
}
`, environmentResourceTFConfig(&environmentResourceTFConfigOptions{TFResourceName: "env", Name: &envName}))
}
