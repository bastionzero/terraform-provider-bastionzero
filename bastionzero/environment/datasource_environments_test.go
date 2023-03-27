package environment_test

import (
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/internal/bztftest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnvironments_Basic(t *testing.T) {
	// Create random env name
	name := bztftest.RandomTestName()

	resourceConfig := environmentResourceTFConfig(&environmentResourceTFConfigOptions{TFResourceName: "env", Name: &name})
	dataSourceConfig := `data "bastionzero_environments" "envs" {}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { bztftest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: bztftest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					// Not much else we can do unless we run against an isolated
					// BastionZero backend. Also, I tried using local variable
					// with filter, but it doesn't seem to be well supported in
					// the terraform-plugin-testing framework.
					bztftest.TestAccCheckListHasElements("data.bastionzero_environments.envs", "environments"),
				),
			},
		},
	})
}
