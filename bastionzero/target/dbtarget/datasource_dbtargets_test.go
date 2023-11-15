package dbtarget_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// whitelist is a list of common attributes between the db_target resource and
// db_target data source schemas
var whitelist = []string{
	"agent_public_key",
	"agent_version",
	"database_authentication_config.cloud_service_provider",
	"database_authentication_config.database",
	"database_authentication_config.authentication_type",
	"database_authentication_config.label",
	"environment_id",
	"id",
	"last_agent_update",
	"name",
	"proxy_target_id",
	"proxy_environment_id",
	"region",
	"remote_host",
	"remote_port",
	"local_port",
	"status",
	"type",
}

func TestAccDbTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	dataSourceName := "data.bastionzero_db_targets.test"

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
			// First create resource
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"),
			},
			// Then check that the list data source contains the db target we
			// created above
			{
				Config: acctest.ConfigCompose(testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"), testAccDbTargetsDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					acctest.CheckListOrSetHasElements(dataSourceName, "targets"),
					// Must whitelist specific attributes because the data
					// source and the resource have slightly different schemas
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName, whitelist, dataSourceName, "targets.*"),
				),
			},
		},
	})
}

func TestAccDbTargetsDataSource_Many(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	dataSourceName := "data.bastionzero_db_targets.test"

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
			// First create many resources
			{
				Config: testAccDbTargetConfigMany(rName, 2, env.ID, bzeroTarget.ID, "localhost", "5432"),
			},
			// Then check that the list data source contains the db targets we
			// created above
			{
				Config: acctest.ConfigCompose(testAccDbTargetConfigMany(rName, 2, env.ID, bzeroTarget.ID, "localhost", "5432"), testAccDbTargetsDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName+".0", new(targets.DatabaseTarget)),
					testAccCheckDbTargetExists(resourceName+".1", new(targets.DatabaseTarget)),
					acctest.CheckListOrSetHasElements(dataSourceName, "targets"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".0", whitelist, dataSourceName, "targets.*"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".1", whitelist, dataSourceName, "targets.*"),
				),
			},
		},
	})
}

func testAccDbTargetsDataSourceConfig() string {
	return `
data "bastionzero_db_targets" "test" {
}
`
}

func testAccDbTargetConfigMany(rName string, count int, envID string, proxyTargetID string, remoteHost string, remotePort string) string {
	return fmt.Sprintf(`
resource "bastionzero_db_target" "test" {
  count = %[2]v	
  environment_id = %[3]q
  name = %[1]q
  proxy_target_id = %[4]q
  remote_host = %[5]q
  remote_port = %[6]q
}
`, rName+"-${count.index}", count, envID, proxyTargetID, remoteHost, remotePort)
}
