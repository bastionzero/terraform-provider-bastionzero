package dbtarget_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSupportedDatabaseConfigsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_supported_database_configs.test"

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSupportedDatabaseConfigsDataSourceConfig(),
				// Basic acceptance test that just checks the data source has
				// some number of values after calling BastionZero and
				// converting the API type to the TF type
				Check: acctest.CheckListOrSetHasElements(dataSourceName, "configs"),
			},
		},
	})
}

func testAccSupportedDatabaseConfigsDataSourceConfig() string {
	return `
data "bastionzero_supported_database_configs" "test" {
}
`
}
