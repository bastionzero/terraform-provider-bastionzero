package targetconnect_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTargetConnectPolicyDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	dataSourceName := "data.bastionzero_targetconnect_policy.test"
	var policy policies.TargetConnectPolicy

	resourcePolicy := testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Shell)})
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: resourcePolicy,
			},
			// Then, check data source matches policy we create
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccTargetConnectPolicyDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "subjects", dataSourceName, "subjects"),
					resource.TestCheckResourceAttrPair(resourceName, "groups", dataSourceName, "groups"),
					resource.TestCheckResourceAttrPair(resourceName, "environments", dataSourceName, "environments"),
					resource.TestCheckResourceAttrPair(resourceName, "targets", dataSourceName, "targets"),
					resource.TestCheckResourceAttrPair(resourceName, "target_users", dataSourceName, "target_users"),
					resource.TestCheckResourceAttrPair(resourceName, "verbs", dataSourceName, "verbs"),
				),
			},
		},
	})
}

func TestTargetConnectPolicyDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccTargetConnectPolicyDataSourceConfigWithID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccTargetConnectPolicyDataSourceConfigWithID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccTargetConnectPolicyDataSourceConfigID() string {
	return `
data "bastionzero_targetconnect_policy" "test" {
  id = bastionzero_targetconnect_policy.test.id
}
`
}

func testAccTargetConnectPolicyDataSourceConfigWithID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_targetconnect_policy" "test" {
  id = %[1]q
}
`, id)
}
