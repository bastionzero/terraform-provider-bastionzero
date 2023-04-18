package sessionrecording_test

import (
	"context"
	"fmt"
	"regexp"
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

func TestSessionRecordingPolicyDataSource_InvalidID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty id not permitted
				Config:      testAccSessionRecordingPolicyDataSourceConfigWithID(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Bad id not permitted
				Config:      testAccSessionRecordingPolicyDataSourceConfigWithID("foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
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

func testAccSessionRecordingPolicyDataSourceConfigWithID(id string) string {
	return fmt.Sprintf(`
data "bastionzero_sessionrecording_policy" "test" {
  id = %[1]q
}
`, id)
}
