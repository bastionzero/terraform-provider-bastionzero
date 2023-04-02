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
			// Then check that the list data source contains the environment we
			// created above
			{
				Config: acctest.ConfigCompose(testAccEnvironmentConfigName(rName), testAccEnvironmentsDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env),
					acctest.CheckListOrSetHasElements(dataSourceName, "environments"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName, []string{}, dataSourceName, "environments.*"),
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
