package environment_test

import (
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bztftest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnvironment_BasicById(t *testing.T) {
	var env environments.Environment
	// Create random env name
	name := bztftest.RandomTestName()

	resourceConfig := environmentResourceTFConfig("env", &environmentResourceOptions{Name: &name})
	dataSourceConfig := `
	data "bastionzero_environment" "env" {
		id = bastionzero_environment.env.id
	}`

	dataSourceRefName := "data.bastionzero_environment.env"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { bztftest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: bztftest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First ensure we can create the resource
			{
				Config: resourceConfig,
			},
			// Read testing
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					bztftest.TestAccCheckEnvironmentExists(dataSourceRefName, &env),
					resource.TestCheckResourceAttr(dataSourceRefName, "name", name),
					resource.TestCheckResourceAttr(dataSourceRefName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceRefName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
					resource.TestCheckResourceAttr(dataSourceRefName, "is_default", "false"),
					resource.TestCheckResourceAttr(dataSourceRefName, "targets.%", "0"),
					resource.TestMatchResourceAttr(dataSourceRefName, "id", bztftest.ExpectedIDRegEx()),
					resource.TestMatchResourceAttr(dataSourceRefName, "organization_id", bztftest.ExpectedIDRegEx()),
					resource.TestMatchResourceAttr(dataSourceRefName, "time_created", bztftest.ExpectedTimestampRegEx()),
				),
			},
		},
	})
}
