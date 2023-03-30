package environment_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_environment.test"
	dataSourceName := "data.bastionzero_environment.test"
	var env environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Check data source matches environment we create
			{
				Config: testAccEnvironmentDataSourceConfigID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "offline_cleanup_timeout_hours", dataSourceName, "offline_cleanup_timeout_hours"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "organization_id", dataSourceName, "organization_id"),
					resource.TestCheckResourceAttrPair(resourceName, "time_created", dataSourceName, "time_created"),
					resource.TestCheckResourceAttrPair(resourceName, "is_default", dataSourceName, "is_default"),
					resource.TestCheckResourceAttrPair(resourceName, "targets", dataSourceName, "targets"),
				),
			},
		},
	})
}

func testAccEnvironmentDataSourceConfigID(rName string) string {
	return acctest.ConfigCompose(
		testAccEnvironmentConfigName(rName),
		`
data "bastionzero_environment" "test" {
  id = bastionzero_environment.test.id
}
`)
}
