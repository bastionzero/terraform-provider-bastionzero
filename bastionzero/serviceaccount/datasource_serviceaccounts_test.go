package serviceaccount_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceAccountsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_service_accounts.test"
	serviceAccount := new(serviceaccounts.ServiceAccount)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNServiceAccountsOrSkip(t, serviceAccount)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountsDataSourceConfig(),
				// Check that the service account we queried for is returned in
				// the list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "service_accounts.*", map[string]string{
					"created_by":       serviceAccount.CreatedBy,
					"email":            serviceAccount.Email,
					"enabled":          fmt.Sprintf("%t", serviceAccount.Enabled),
					"external_id":      serviceAccount.ExternalID,
					"id":               serviceAccount.ID,
					"is_admin":         fmt.Sprintf("%t", serviceAccount.IsAdmin),
					"jwks_url":         serviceAccount.JwksURL,
					"jwks_url_pattern": serviceAccount.JwksURLPattern,
					"last_login":       serviceAccount.LastLogin.UTC().Format(time.RFC3339),
					"organization_id":  serviceAccount.OrganizationID,
					"time_created":     serviceAccount.TimeCreated.UTC().Format(time.RFC3339),
					"type":             string(subjecttype.ServiceAccount),
				}),
			},
		},
	})
}

func testAccServiceAccountsDataSourceConfig() string {
	return `
data "bastionzero_service_accounts" "test" {
}
`
}
