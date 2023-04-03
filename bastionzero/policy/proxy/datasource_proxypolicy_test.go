package proxy_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProxyPolicyDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	dataSourceName := "data.bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy

	resourcePolicy := testAccProxyPolicyConfigBasic(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: resourcePolicy,
			},
			// Then, check data source matches policy we create
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccProxyPolicyDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "subjects", dataSourceName, "subjects"),
					resource.TestCheckResourceAttrPair(resourceName, "groups", dataSourceName, "groups"),
					resource.TestCheckResourceAttrPair(resourceName, "environments", dataSourceName, "environments"),
					resource.TestCheckResourceAttrPair(resourceName, "targets", dataSourceName, "targets"),
					resource.TestCheckResourceAttrPair(resourceName, "target_users", dataSourceName, "target_users"),
				),
			},
		},
	})
}

func testAccProxyPolicyDataSourceConfigID() string {
	return `
data "bastionzero_proxy_policy" "test" {
  id = bastionzero_proxy_policy.test.id
}
`
}
