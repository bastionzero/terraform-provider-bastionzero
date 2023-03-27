package environment_test

import (
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/internal/bztftest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnvironments_Basic(t *testing.T) {
	// Create random env name
	name := bztftest.RandomTestName()

	resourceConfig := environmentResourceTFConfig("env", &environmentResourceOptions{Name: &name})
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
					//
					// We can probably test more attributes if we force TF
					// version 1.4.0 and use `terraform_data` resource (coupled
					// with local variable that filters for env with name), but
					// don't want to add TF specific tests until this issue is
					// resolved:
					// https://github.com/hashicorp/terraform-plugin-testing/issues/68
					bztftest.TestAccCheckListHasElements("data.bastionzero_environments.envs", "environments"),
				),
			},
		},
	})
}
