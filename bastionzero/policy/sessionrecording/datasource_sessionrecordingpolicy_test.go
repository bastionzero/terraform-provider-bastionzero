package sessionrecording_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSessionRecordingPolicyDataSource_ID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	dataSourceName := "data.bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy

	resourcePolicy := testAccSessionRecordingPolicyConfigBasic(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// First create a resource
			{
				Config: resourcePolicy,
			},
			// Then, check data source matches policy we create
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccSessionRecordingPolicyDataSourceConfigID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "subjects", dataSourceName, "subjects"),
					resource.TestCheckResourceAttrPair(resourceName, "groups", dataSourceName, "groups"),
					resource.TestCheckResourceAttrPair(resourceName, "record_input", dataSourceName, "record_input"),
				),
			},
		},
	})
}

func testAccSessionRecordingPolicyDataSourceConfigID() string {
	return `
data "bastionzero_sessionrecording_policy" "test" {
  id = bastionzero_sessionrecording_policy.test.id
}
`
}