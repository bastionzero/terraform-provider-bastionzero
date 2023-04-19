package environment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
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

func TestAccEnvironmentsDataSource_Many(t *testing.T) {
	ctx := context.Background()
	resourceName := "bastionzero_environment.test"
	dataSourceName := "data.bastionzero_environments.test"
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			// First create many resources
			{
				Config: testAccEnvironmentConfigMany(rName, 2),
			},
			// Then check that the list data source contains the environments we
			// created above
			{
				Config: acctest.ConfigCompose(testAccEnvironmentConfigMany(rName, 2), testAccEnvironmentsDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName+".0", new(environments.Environment)),
					testAccCheckEnvironmentExists(resourceName+".1", new(environments.Environment)),
					acctest.CheckListOrSetHasElements(dataSourceName, "environments"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".0", []string{}, dataSourceName, "environments.*"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".1", []string{}, dataSourceName, "environments.*"),
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

func testAccEnvironmentConfigMany(rName string, count int) string {
	return fmt.Sprintf(`
resource "bastionzero_environment" "test" {
  count = %[2]v
  name = %[1]q
}
`, rName+"-${count.index}", count)
}
