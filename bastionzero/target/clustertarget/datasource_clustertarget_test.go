package clustertarget_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCheckValidClusterUsers(resourceName string, validClusterUsers []string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(resourceName, "valid_cluster_users.#", strconv.Itoa(len(validClusterUsers)))}
	for _, validUser := range validClusterUsers {
		checks = append(checks, resource.TestCheckTypeSetElemAttr(
			resourceName,
			"valid_cluster_users.*",
			validUser,
		))
	}

	return resource.ComposeTestCheckFunc(checks...)
}

func TestAccClusterTargetDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_cluster_target.test"
	clusterTarget := new(targets.ClusterTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNClusterTargetsOrSkip(t, clusterTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTargetDataSourceConfigID(clusterTarget.ID),
				// Check the data source attributes look correct based on the
				// Cluster target we queried for
				Check: resource.ComposeTestCheckFunc(
					acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, clusterTarget, getValuesCheckMap),
					testAccCheckValidClusterUsers(dataSourceName, clusterTarget.ValidClusterUsers),
				),
			},
		},
	})
}

func testAccClusterTargetDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_cluster_target" "test" {
  id = %[1]q
}
`, id)
}
