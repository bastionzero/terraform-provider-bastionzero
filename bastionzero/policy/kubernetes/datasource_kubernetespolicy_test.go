package kubernetes_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKubernetesPolicyDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	dataSourceName := "data.bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resourcePolicy := testAccKubernetesPolicyConfigBasic(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: resourcePolicy,
			},
			// Then, check data source matches policy we create
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccKubernetesPolicyDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "subjects", dataSourceName, "subjects"),
					resource.TestCheckResourceAttrPair(resourceName, "groups", dataSourceName, "groups"),
					resource.TestCheckResourceAttrPair(resourceName, "environments", dataSourceName, "environments"),
					resource.TestCheckResourceAttrPair(resourceName, "clusters", dataSourceName, "clusters"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_users", dataSourceName, "cluster_users"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_groups", dataSourceName, "cluster_groups"),
				),
			},
		},
	})
}

func TestKubernetesPolicyDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccKubernetesPolicyDataSourceConfigWithID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccKubernetesPolicyDataSourceConfigWithID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccKubernetesPolicyDataSourceConfigID() string {
	return `
data "bastionzero_kubernetes_policy" "test" {
  id = bastionzero_kubernetes_policy.test.id
}
`
}

func testAccKubernetesPolicyDataSourceConfigWithID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_kubernetes_policy" "test" {
  id = %[1]q
}
`, id)
}
