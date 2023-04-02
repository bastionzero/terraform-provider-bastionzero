package targetconnect_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceTargetConnectPolicies_Basic(t *testing.T) {
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

func testAccTargetConnectPoliciesDataSourceConfig() string {
	return `
data "bastionzero_targetconnect_policies" "test" {
}
`
}

func TestAccDataSourceTargetConnectPolicies_FilterSubjects(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_targetconnect_policies.test"
	subject := new(policies.Subject)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNUsersOrSkip(t, subject)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(testAccTargetConnectPoliciesDataSourceConfigFilterSubjects([]string{subject.ID})),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckListOrSetHasElements(dataSourceName, "policies"),
					CheckAllPoliciesHaveSubjectID(dataSourceName, subject.ID),
				),
			},
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

// CheckAllPoliciesHaveSubjectID checks that all policies have at least one
// subject that matches an expected ID
func CheckAllPoliciesHaveSubjectID(namedTFResource, expectedSubjectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		totalPolicies, err := acctest.ListOrSetCount(rs, "policies")
		if err != nil {
			return err
		}

		if totalPolicies == 0 {
			return fmt.Errorf("list of policies is empty")
		}

		// Aggregate attribute checked errors
		var result *multierror.Error

	POLICY:
		for i := 0; i < totalPolicies; i++ {
			totalSubjects, err := acctest.ListOrSetCount(rs, fmt.Sprintf("policies.%v.subjects", i))
			if err != nil {
				return err
			}

			for j := 0; j < totalSubjects; j++ {
				if err := resource.TestCheckResourceAttr(namedTFResource, fmt.Sprintf("policies.%v.subjects.%v.id", i, j), expectedSubjectID)(s); err == nil {
					// Found at least one! Continue checking the next policy
					continue POLICY
				} else {
					// This subject does not match. Aggregate this error.
					result = multierror.Append(result, err)
				}
			}
		}

		return result.ErrorOrNil()
	}
}
