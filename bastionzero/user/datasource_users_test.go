package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(user *users.User) map[string]string {
	valuesCheckMap := map[string]string{
		"email":           user.Email,
		"full_name":       user.FullName,
		"id":              user.ID,
		"is_admin":        fmt.Sprintf("%t", user.IsAdmin),
		"organization_id": user.OrganizationID,
		"time_created":    user.TimeCreated.UTC().Format(time.RFC3339),
		"type":            string(subjecttype.ServiceAccount),
	}

	if user.LastLogin != nil {
		valuesCheckMap["last_login"] = user.LastLogin.UTC().Format(time.RFC3339)
	} else {
		// TestCheckTypeSetElemNestedAttrs() will check to see that this
		// attribute is unset (null)
		valuesCheckMap["last_login"] = ""
	}

	return valuesCheckMap
}

func TestAccUsersDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_users.test"
	user := new(users.User)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkip(t, user)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig(),
				// Check that the user we queried for is returned in the list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "users.*", getValuesCheckMap(user)),
			},
		},
	})
}

func testAccUsersDataSourceConfig() string {
	return `
data "bastionzero_users" "test" {
}
`
}
