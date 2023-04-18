package dbtarget_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDbTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_db_target.test"
	dbTarget := new(targets.DatabaseTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNDbTargetsOrSkip(t, dbTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetDataSourceConfigID(dbTarget.ID),
				// Check the data source attributes look correct based on the Db
				// target we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, dbTarget, getValuesCheckMap),
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
				Config:      testAccDbTargetDataSourceConfigID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccDbTargetDataSourceConfigID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccDbTargetDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_db_target" "test" {
  id = %[1]q
}
`, id)
}
