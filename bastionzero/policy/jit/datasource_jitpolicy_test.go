package jit_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJITPolicyDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	dataSourceName := "data.bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)

	resourcePolicy := testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID})
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: resourcePolicy,
			},
			// Then, check data source matches policy we create
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccJITPolicyDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "subjects", dataSourceName, "subjects"),
					resource.TestCheckResourceAttrPair(resourceName, "groups", dataSourceName, "groups"),
					resource.TestCheckResourceAttrPair(resourceName, "auto_approved", dataSourceName, "auto_approved"),
					resource.TestCheckResourceAttrPair(resourceName, "child_policies", dataSourceName, "child_policies"),
					resource.TestCheckResourceAttrPair(resourceName, "duration", dataSourceName, "duration"),
				),
			},
		},
	})
}

func testAccJITPolicyDataSourceConfigID() string {
	return `
data "bastionzero_jit_policy" "test" {
  id = bastionzero_jit_policy.test.id
}
`
}
