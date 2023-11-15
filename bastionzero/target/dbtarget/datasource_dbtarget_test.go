package dbtarget_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	dataSourceName := "data.bastionzero_db_target.test"

	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"),
			},
			// Then, check data source matches db target we create
			{
				Config: acctest.ConfigCompose(testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"), testAccDbTargetDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", dataSourceName, "environment_id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "proxy_target_id", dataSourceName, "proxy_target_id"),
					resource.TestCheckResourceAttrPair(resourceName, "proxy_environment_id", dataSourceName, "proxy_environment_id"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_host", dataSourceName, "remote_host"),
					resource.TestCheckResourceAttrPair(resourceName, "remote_port", dataSourceName, "remote_port"),
					resource.TestCheckResourceAttrPair(resourceName, "database_authentication_config.authentication_type", dataSourceName, "database_authentication_config.authentication_type"),
					resource.TestCheckResourceAttrPair(resourceName, "database_authentication_config.label", dataSourceName, "database_authentication_config.label"),
					resource.TestCheckResourceAttrPair(resourceName, "database_authentication_config.cloud_service_provider", dataSourceName, "database_authentication_config.cloud_service_provider"),
					resource.TestCheckResourceAttrPair(resourceName, "database_authentication_config.database", dataSourceName, "database_authentication_config.database"),
					resource.TestCheckResourceAttrPair(resourceName, "local_port", dataSourceName, "local_port"),
				),
			},
		},
	})
}

func TestDbTargetDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccDbTargetDataSourceConfigWithID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccDbTargetDataSourceConfigWithID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccDbTargetDataSourceConfigID() string {
	return `
data "bastionzero_db_target" "test" {
  id = bastionzero_db_target.test.id
}
`
}

func testAccDbTargetDataSourceConfigWithID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_db_target" "test" {
  id = %[1]q
}
`, id)
}
