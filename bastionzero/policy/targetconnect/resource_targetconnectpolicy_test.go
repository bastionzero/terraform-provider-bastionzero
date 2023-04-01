package targetconnect_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccTargetConnectPolicy_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy policies.TargetConnectPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo", "bar"}, []string{"Shell"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					testAccCheckTargetConnectPolicyAttributes(t, &policy, &expectedTargetConnectPolicy{
						Name:         &rName,
						Description:  bastionzero.PtrTo(""),
						Subjects:     &[]policies.Subject{},
						Groups:       &[]policies.Group{},
						Environments: &[]policies.Environment{},
						Targets:      &[]policies.Target{},
						TargetUsers:  &[]policies.TargetUser{{Username: "foo"}, {Username: "bar"}},
						Verbs:        &[]policies.Verb{{Type: verbtype.Shell}},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "foo"),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "bar"),
					resource.TestCheckTypeSetElemAttr(resourceName, "verbs.*", "Shell"),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "type", string(policytype.TargetConnect)),
					// Check that unspecified values remain null
					resource.TestCheckNoResourceAttr(resourceName, "subjects"),
					resource.TestCheckNoResourceAttr(resourceName, "groups"),
					resource.TestCheckNoResourceAttr(resourceName, "environments"),
					resource.TestCheckNoResourceAttr(resourceName, "targets"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTargetConnectPolicyConfigBasic(rName string, targetUsers []string, verbs []string) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformList(targetUsers), acctest.ToTerraformList(verbs))
}

type expectedTargetConnectPolicy struct {
	Name         *string
	Description  *string
	Subjects     *[]policies.Subject
	Groups       *[]policies.Group
	Environments *[]policies.Environment
	Targets      *[]policies.Target
	TargetUsers  *[]policies.TargetUser
	Verbs        *[]policies.Verb
}

func testAccCheckTargetConnectPolicyAttributes(t *testing.T, policy *policies.TargetConnectPolicy, expected *expectedTargetConnectPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if expected.Name != nil && *expected.Name != policy.Name {
			return fmt.Errorf("Bad name, expected \"%s\", got: %#v", *expected.Name, policy.Name)
		}
		if expected.Description != nil && *expected.Description != policy.GetDescription() {
			return fmt.Errorf("Bad description, expected \"%s\", got: %#v", *expected.Description, policy.Description)
		}
		if expected.Subjects != nil && !assert.ElementsMatch(t, *expected.Subjects, policy.GetSubjects()) {
			return fmt.Errorf("Bad subjects, expected \"%s\", got: %#v", *expected.Subjects, policy.Subjects)
		}
		if expected.Groups != nil && !assert.ElementsMatch(t, *expected.Groups, policy.GetGroups()) {
			return fmt.Errorf("Bad groups, expected \"%s\", got: %#v", *expected.Groups, policy.Groups)
		}
		if expected.Environments != nil && !assert.ElementsMatch(t, *expected.Environments, policy.GetEnvironments()) {
			return fmt.Errorf("Bad environments, expected \"%s\", got: %#v", *expected.Environments, policy.Environments)
		}
		if expected.Targets != nil && !assert.ElementsMatch(t, *expected.Targets, policy.GetTargets()) {
			return fmt.Errorf("Bad targets, expected \"%s\", got: %#v", *expected.Targets, policy.Targets)
		}
		if expected.TargetUsers != nil && !assert.ElementsMatch(t, *expected.TargetUsers, policy.GetTargetUsers()) {
			return fmt.Errorf("Bad target_users, expected \"%s\", got: %#v", *expected.TargetUsers, policy.TargetUsers)
		}
		if expected.Verbs != nil && !assert.ElementsMatch(t, *expected.Verbs, policy.GetVerbs()) {
			return fmt.Errorf("Bad verbs, expected \"%s\", got: %#v", *expected.Verbs, policy.Verbs)
		}

		return nil
	}
}

func testAccCheckTargetConnectPolicyExists(namedTFResource string, policy *policies.TargetConnectPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.TargetConnectPolicy, *http.Response, error) {
		return c.Policies.GetTargetConnectPolicy(ctx, id)
	})
}

func testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
	)
}

func testAccCheckTargetConnectPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bastionzero_targetconnect_policy" {
			continue
		}

		// Try to find the policy
		_, _, err := acctest.APIClient.Policies.GetTargetConnectPolicy(context.Background(), rs.Primary.ID)
		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("Error waiting for target connect policy (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
