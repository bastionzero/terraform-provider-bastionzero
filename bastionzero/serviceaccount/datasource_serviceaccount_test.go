package serviceaccount_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceAccountDataSource_ID(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_service_account.test"
	serviceAccount := new(serviceaccounts.ServiceAccount)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNServiceAccountsOrSkip(t, serviceAccount)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountDataSourceConfigID(serviceAccount.ID),
				// Check the data source attributes look correct based on the
				// service account we queried for
				Check: acctest.ExpandValuesCheckMapToSingleCheck(dataSourceName, serviceAccount, getValuesCheckMap),
			},
		},
	})
}

func TestServiceAccountDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccServiceAccountDataSourceConfigID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccServiceAccountDataSourceConfigID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccServiceAccountDataSourceConfigID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_service_account" "test" {
  id = %[1]q
}
`, id)
}
