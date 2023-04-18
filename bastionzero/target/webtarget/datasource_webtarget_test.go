package webtarget_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_web_target.test"
	webTarget := new(targets.WebTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNWebTargetsOrSkip(t, webTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebTargetDataSourceConfigID(webTarget.ID),
				// Check the data source attributes look correct based on the
				// Web target we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, webTarget, getValuesCheckMap),
			},
		},
	})
}

func TestWebTargetDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccWebTargetDataSourceConfigID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccWebTargetDataSourceConfigID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccWebTargetDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_web_target" "test" {
  id = %[1]q
}
`, id)
}
