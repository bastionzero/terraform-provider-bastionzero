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

	// schema := environment.MakeEnvironmentResourceSchema()
	// keys := make([]string, 0, len(schema))
	// for k := range schema {
	// 	keys = append(keys, k)
	// }

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
					acctest.CheckListHasElements(dataSourceName, "environments"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(t, resourceName, []string{}, dataSourceName, "environments.*"),
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
