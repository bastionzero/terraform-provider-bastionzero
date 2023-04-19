package targetconnect_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTargetConnectPoliciesDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	dataSourceName := "data.bastionzero_targetconnect_policies.test"
	var policy policies.TargetConnectPolicy

	resourcePolicy := testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Shell)})
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that the list data source contains the policy we
			// created above
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccTargetConnectPoliciesDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName, []string{}, dataSourceName, "policies.*"),
				),
			},
		},
	})
}

func TestAccTargetConnectPoliciesDataSource_Many(t *testing.T) {
	ctx := context.Background()
	resourceName := "bastionzero_targetconnect_policy.test"
	dataSourceName := "data.bastionzero_targetconnect_policies.test"
	rName := acctest.RandomName()

	resourcePolicy := testAccTargetConnectPolicyConfigMany(rName, []string{"foo"}, []string{string(verbtype.Shell)}, 2)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// First create many resources
			{
				Config: resourcePolicy,
			},
			// Then check that the list data source contains the policies we
			// created above
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccTargetConnectPoliciesDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName+".0", new(policies.TargetConnectPolicy)),
					testAccCheckTargetConnectPolicyExists(resourceName+".1", new(policies.TargetConnectPolicy)),
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".0", []string{}, dataSourceName, "policies.*"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName+".1", []string{}, dataSourceName, "policies.*"),
				),
			},
		},
	})
}

func testAccTargetConnectPoliciesDataSourceConfig() string {
	return `
data "bastionzero_targetconnect_policies" "test" {
}
`
}

func testAccTargetConnectPolicyConfigMany(rName string, targetUsers []string, verbs []string, count int) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  count = %[4]v
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName+"-${count.index}", acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), count)
}

func TestAccTargetConnectPoliciesDataSource_FilterSubjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_targetconnect_policies.test"
	subject := new(policies.Subject)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkipAsPolicySubject(t, subject)

	resourcePolicy := testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccTargetConnectPoliciesDataSourceConfigFilterSubjects([]string{subject.ID})),
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

func testAccTargetConnectPoliciesDataSourceConfigFilterSubjects(subjectIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_targetconnect_policies" "test" {
  filter_subjects = %[1]s
}
`, acctest.ToTerraformStringList(subjectIDs))
}

func TestAccTargetConnectPoliciesDataSource_FilterGroups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_targetconnect_policies.test"
	group := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group)

	resourcePolicy := testAccTargetConnectPolicyConfigGroups(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyGroups(ctx, []policies.Group{*group}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccTargetConnectPoliciesDataSourceConfigFilterGroups([]string{group.ID})),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckAllPoliciesHaveGroupID(dataSourceName, group.ID),
				),
			},
			// Zero matches
			{
				Config: testAccTargetConnectPoliciesDataSourceConfigFilterGroups([]string{"foo"}),
				Check:  resource.TestCheckResourceAttr(dataSourceName, "policies.#", "0"),
			},
		},
	})
}

func testAccTargetConnectPoliciesDataSourceConfigFilterGroups(groupIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_targetconnect_policies" "test" {
  filter_groups = %[1]s
}
`, acctest.ToTerraformStringList(groupIDs))
}
