package sessionrecording_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSessionRecordingPoliciesDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	dataSourceName := "data.bastionzero_sessionrecording_policies.test"
	var policy policies.SessionRecordingPolicy

	resourcePolicy := testAccSessionRecordingPolicyConfigBasic(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that the list data source contains the policy we
			// created above
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccSessionRecordingPoliciesDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName, []string{}, dataSourceName, "policies.*"),
				),
			},
		},
	})
}

func TestAccSessionRecordingPoliciesDataSource_Many(t *testing.T) {
	ctx := context.Background()
	resourceName := "bastionzero_sessionrecording_policy.test"
	dataSourceName := "data.bastionzero_sessionrecording_policies.test"
	rName := acctest.RandomName()

	resourcePolicy := testAccSessionRecordingPolicyConfigMany(rName, 2)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// First create many resources
			{
				Config: resourcePolicy,
			},
			// Then check that the list data source contains the policies we
			// created above
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccSessionRecordingPoliciesDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName+".0", new(policies.SessionRecordingPolicy)),
					testAccCheckSessionRecordingPolicyExists(resourceName+".1", new(policies.SessionRecordingPolicy)),
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".0", []string{}, dataSourceName, "policies.*"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".1", []string{}, dataSourceName, "policies.*"),
				),
			},
		},
	})
}

func testAccSessionRecordingPoliciesDataSourceConfig() string {
	return `
data "bastionzero_sessionrecording_policies" "test" {
}
`
}

func testAccSessionRecordingPolicyConfigMany(rName string, count int) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  count = %[2]v
  name = %[1]q
}
`, rName+"-${count.index}", count)
}

func TestAccSessionRecordingPoliciesDataSource_FilterSubjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_sessionrecording_policies.test"
	subject := new(policies.Subject)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkipAsPolicySubject(t, subject)

	resourcePolicy := testAccSessionRecordingPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccSessionRecordingPoliciesDataSourceConfigFilterSubjects([]string{subject.ID})),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckAllPoliciesHaveSubjectID(dataSourceName, subject.ID),
				),
			},
			// Cannot do zero matches test because must provide valid subject
			// UUID which we can't guarantee. Can add later if we remove backend
			// restriction that subject ID must exist
		},
	})
}

func testAccSessionRecordingPoliciesDataSourceConfigFilterSubjects(subjectIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_sessionrecording_policies" "test" {
  filter_subjects = %[1]s
}
`, acctest.ToTerraformStringList(subjectIDs))
}

func TestAccSessionRecordingPoliciesDataSource_FilterGroups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_sessionrecording_policies.test"
	group := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group)

	resourcePolicy := testAccSessionRecordingPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccSessionRecordingPoliciesDataSourceConfigFilterGroups([]string{group.ID})),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckAllPoliciesHaveGroupID(dataSourceName, group.ID),
				),
			},
			// Zero matches
			{
				Config: testAccSessionRecordingPoliciesDataSourceConfigFilterGroups([]string{"foo"}),
				Check:  resource.TestCheckResourceAttr(dataSourceName, "policies.#", "0"),
			},
		},
	})
}

func testAccSessionRecordingPoliciesDataSourceConfigFilterGroups(groupIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_sessionrecording_policies" "test" {
  filter_groups = %[1]s
}
`, acctest.ToTerraformStringList(groupIDs))
}
