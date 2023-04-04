package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_user.test"
	user := new(users.User)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkip(t, user)

	t.Fail()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfigID(user.ID),
				// Check the data source attributes look correct based on the
				// user we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, user, getValuesCheckMap),
			},
		},
	})
}

func testAccUserDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_user" "test" {
  id = %[1]q
}
`, id)
}
