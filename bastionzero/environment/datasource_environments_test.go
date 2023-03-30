package environment_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceEnvironments_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_environment.test"
	dataSourceName := "data.bastionzero_environments.test"
	var env environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: testAccEnvironmentConfigName(rName),
			},
			{
				Config: acctest.ConfigCompose(testAccEnvironmentConfigName(rName), testAccEnvironmentsDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env),
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
					acctest.CheckListHasElements(dataSourceName, "environments"),
					resource.TestCheckTypeSetElemAttr(dataSourceName, "environments.*", "value1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "name", dataSourceName, "environments.*"),
				),
			},
		},
	})
}

func testAccEnvironmentsDataSourceConfig() string {
	return `
data "bastionzero_environments" "test" {
}
`
}
