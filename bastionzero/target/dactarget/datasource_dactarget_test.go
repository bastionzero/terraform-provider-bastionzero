package dactarget_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDACTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_dac_target.test"
	dacTarget := new(targets.DynamicAccessConfiguration)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNDACTargetsOrSkip(t, dacTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDACTargetDataSourceConfigID(dacTarget.ID),
				// Check the data source attributes look correct based on the
				// DAC target we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, dacTarget, getValuesCheckMap),
			},
		},
	})
}

func testAccDACTargetDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_dac_target" "test" {
  id = %[1]q
}
`, id)
}
