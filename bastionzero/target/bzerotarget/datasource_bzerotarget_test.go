package bzerotarget_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBzeroTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_bzero_target.test"
	bzeroTarget := new(targets.BzeroTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBzeroTargetDataSourceConfigID(bzeroTarget.ID),
				// Check the data source attributes look correct based on the
				// Bzero target we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, bzeroTarget, getValuesCheckMap),
			},
		},
	})
}

func testAccBzeroTargetDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_bzero_target" "test" {
  id = %[1]q
}
`, id)
}
