package organization_test

import (
	"context"
	"testing"

	bzapi "github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/organization"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_groups.test"
	group := new(bzapi.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkip(t, group)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSourceConfig(),
				// Check that the group we queried for is returned in the list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "groups.*", map[string]string{"id": group.ID, "name": group.Name}),
			},
		},
	})
}

func testAccGroupsDataSourceConfig() string {
	return `
data "bastionzero_groups" "test" {
}
`
}
