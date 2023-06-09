package environment_test

import (
	"context"
	"fmt"
	"regexp"
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
			// First create a resource
			{
				Config: testAccEnvironmentConfigName(rName),
			},
			// Then, check data source matches environment we create
			{
				Config: acctest.ConfigCompose(testAccEnvironmentConfigName(rName), testAccEnvironmentDataSourceConfigID()),
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

func TestEnvironmentDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccEnvironmentDataSourceConfigWithID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccEnvironmentDataSourceConfigWithID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccEnvironmentDataSourceConfigID() string {
	return `
data "bastionzero_environment" "test" {
  id = bastionzero_environment.test.id
}
`
}

func testAccEnvironmentDataSourceConfigWithID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_environment" "test" {
  id = %[1]q
}
`, id)
}
