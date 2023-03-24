package environment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
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
					resource.TestCheckResourceAttr("data.bastionzero_environment.env", "name", name),
					resource.TestCheckResourceAttr("data.bastionzero_environment.env", "description", ""),
					resource.TestCheckResourceAttr("data.bastionzero_environment.env", "offline_cleanup_timeout_hours", fmt.Sprint(90*24*time.Hour.Hours())),
					resource.TestCheckResourceAttr("data.bastionzero_environment.env", "is_default", "false"),
					resource.TestCheckResourceAttr("data.bastionzero_environment.env", "targets.%", "0"),
					resource.TestMatchResourceAttr("data.bastionzero_environment.env", "id", bztftest.ExpectedIDRegEx()),
					resource.TestMatchResourceAttr("data.bastionzero_environment.env", "organization_id", bztftest.ExpectedIDRegEx()),
					resource.TestMatchResourceAttr("data.bastionzero_environment.env", "time_created", bztftest.ExpectedTimestampRegEx()),
				),
			},
		},
	})
}

func environmentByIdConfig(envName string) string {
	return fmt.Sprintf(`
resource "bastionzero_environment" "env" {
  name   = "%s"
}
data "bastionzero_environment" "env" {
  id = bastionzero_environment.env.id
}
`, envName)
}
