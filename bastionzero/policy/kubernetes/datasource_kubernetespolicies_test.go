package kubernetes_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKubernetesPoliciesDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	dataSourceName := "data.bastionzero_kubernetes_policies.test"
	var policy policies.KubernetesPolicy

	resourcePolicy := testAccKubernetesPolicyConfigBasic(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that the list data source contains the policy we
			// created above
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccKubernetesPoliciesDataSourceConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckTypeSetElemNestedAttrsFromResource(resourceName, []string{}, dataSourceName, "policies.*"),
				),
			},
		},
	})
}

func testAccKubernetesPoliciesDataSourceConfig() string {
	return `
data "bastionzero_kubernetes_policies" "test" {
}
`
}

func TestAccKubernetesPoliciesDataSource_FilterSubjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_kubernetes_policies.test"
	subject := new(policies.Subject)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkip(t, subject)

	resourcePolicy := testAccKubernetesPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccKubernetesPoliciesDataSourceConfigFilterSubjects([]string{subject.ID})),
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

func testAccKubernetesPoliciesDataSourceConfigFilterSubjects(subjectIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_kubernetes_policies" "test" {
  filter_subjects = %[1]s
}
`, acctest.ToTerraformStringList(subjectIDs))
}

func TestAccKubernetesPoliciesDataSource_FilterGroups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	dataSourceName := "data.bastionzero_kubernetes_policies.test"
	group := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group)

	resourcePolicy := testAccKubernetesPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group}))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			// First create resource
			{
				Config: resourcePolicy,
			},
			// Then check that we can filter for it
			{
				Config: acctest.ConfigCompose(resourcePolicy, testAccKubernetesPoliciesDataSourceConfigFilterGroups([]string{group.ID})),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					acctest.CheckAllPoliciesHaveGroupID(dataSourceName, group.ID),
				),
			},
			// Zero matches
			{
				Config: testAccKubernetesPoliciesDataSourceConfigFilterGroups([]string{"foo"}),
				Check:  resource.TestCheckResourceAttr(dataSourceName, "policies.#", "0"),
			},
		},
	})
}

func testAccKubernetesPoliciesDataSourceConfigFilterGroups(groupIDs []string) string {
	return fmt.Sprintf(`
data "bastionzero_kubernetes_policies" "test" {
  filter_groups = %[1]s
}
`, acctest.ToTerraformStringList(groupIDs))
}
