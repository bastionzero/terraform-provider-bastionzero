package jit_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func getChildPoliciesOrSkip(ctx context.Context, t *testing.T) (tcPolicy *policies.TargetConnectPolicy, kubePolicy *policies.KubernetesPolicy, proxyPolicy *policies.ProxyPolicy) {
	// We could create the policies using the respective resource types, but I
	// want this test suite to work in isolation of whether those resources work
	tcPolicy = new(policies.TargetConnectPolicy)
	kubePolicy = new(policies.KubernetesPolicy)
	proxyPolicy = new(policies.ProxyPolicy)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNTargetConnectPoliciesOrSkip(t, tcPolicy)
	acctest.FindNKubernetesPoliciesOrSkip(t, kubePolicy)
	acctest.FindNProxyPoliciesOrSkip(t, proxyPolicy)
	return
}

func convertToChildPolicies(ps ...policies.PolicyInterface) []policies.ChildPolicy {
	childPolicies := make([]policies.ChildPolicy, 0)
	for _, policy := range ps {
		childPolicies = append(childPolicies, policies.ChildPolicy{ID: policy.GetID(), Type: policy.GetPolicyType(), Name: policy.GetName()})
	}
	return childPolicies
}

func testAccCheckResourceJITPolicyChildPolicies(resourceName string, childPolicies ...policies.ChildPolicy) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(resourceName, "child_policies.#", strconv.Itoa(len(childPolicies)))}
	for _, childPolicy := range childPolicies {
		checks = append(checks, resource.TestCheckTypeSetElemNestedAttrs(
			resourceName,
			"groups.*",
			map[string]string{
				"id":   childPolicy.ID,
				"name": childPolicy.Name,
				"type": string(childPolicy.Type),
			}),
		)
	}

	return resource.ComposeTestCheckFunc(checks...)
}

func toChildPoliciesSet(childPolicyIDs []string) types.Set {
	// Anonymous type with just the required attributes for child_policies
	// attribute
	type requiredChildPolicyModel struct {
		ID types.String `tfsdk:"id"`
	}
	attributeTypes, _ := internal.AttributeTypes[requiredChildPolicyModel](context.Background())
	elementType := types.ObjectType{AttrTypes: attributeTypes}

	return internal.FlattenFrameworkSet(context.Background(), elementType, childPolicyIDs, func(id string) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id": types.StringValue(id),
		})
	})
}

func TestAccJITPolicy_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)
	asChildPolicies := convertToChildPolicies(tcPolicy, kubePolicy, proxyPolicy)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					testAccCheckTargetConnectPolicyAttributes(t, &policy, &expectedTargetConnectPolicy{
						Name:                  &rName,
						Description:           bastionzero.PtrTo(""),
						Subjects:              &[]policies.Subject{},
						Groups:                &[]policies.Group{},
						ChildPolicies:         &asChildPolicies,
						AutomaticallyApproved: bastionzero.PtrTo(false),
						Duration:              bastionzero.PtrTo(uint(60)),
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckResourceJITPolicyChildPolicies(resourceName, asChildPolicies...),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_approved", "false"),
					resource.TestCheckResourceAttr(resourceName, "duration", "60"),
					// Check that unspecified values remain null
					resource.TestCheckNoResourceAttr(resourceName, "subjects"),
					resource.TestCheckNoResourceAttr(resourceName, "groups"),
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

func testAccTargetConnectPolicyConfigBasic(rName string, childPolicyIDs []string) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String())
}

func testAccTargetConnectPolicyConfigDescription(rName string, childPolicyIDs []string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  description = %[3]q
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), description)
}

func testAccTargetConnectPolicyConfigSubjects(rName string, childPolicyIDs []string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  subjects = %[3]s
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), subjects.String())
}

func testAccTargetConnectPolicyConfigGroups(rName string, childPolicyIDs []string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  groups = %[3]s
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), groups.String())
}

func testAccJITPolicyConfigAutoApproved(rName string, childPolicyIDs []string, autoApproved bool) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  auto_approved = %[3]t
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), autoApproved)
}

func testAccJITPolicyConfigDuration(rName string, childPolicyIDs []string, duration uint) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  duration = %[3]v
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), duration)
}

type expectedTargetConnectPolicy struct {
	Name                  *string
	Description           *string
	Subjects              *[]policies.Subject
	Groups                *[]policies.Group
	ChildPolicies         *[]policies.ChildPolicy
	AutomaticallyApproved *bool
	Duration              *uint
}

func testAccCheckTargetConnectPolicyAttributes(t *testing.T, policy *policies.JITPolicy, expected *expectedTargetConnectPolicy) resource.TestCheckFunc {
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
		if expected.ChildPolicies != nil && !assert.ElementsMatch(t, *expected.ChildPolicies, policy.GetChildPolicies()) {
			return fmt.Errorf("Bad child_policies, expected \"%s\", got: %#v", *expected.ChildPolicies, policy.ChildPolicies)
		}
		if expected.AutomaticallyApproved != nil && *expected.AutomaticallyApproved != policy.GetAutomaticallyApproved() {
			return fmt.Errorf("Bad auto_approved, expected \"%t\", got: %#v", *expected.AutomaticallyApproved, policy.AutomaticallyApproved)
		}
		if expected.Duration != nil && *expected.Duration != policy.GetDuration() {
			return fmt.Errorf("Bad duration, expected \"%d\", got: %#v", *expected.Duration, policy.Duration)
		}

		return nil
	}
}

func testAccCheckTargetConnectPolicyExists(namedTFResource string, policy *policies.JITPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.JITPolicy, *http.Response, error) {
		return c.Policies.GetJITPolicy(ctx, id)
	})
}

func testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.JustInTime)),
	)
}

func testAccCheckTargetConnectPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bastionzero_jit_policy" {
			continue
		}

		// Try to find the policy
		_, _, err := acctest.APIClient.Policies.GetJITPolicy(context.Background(), rs.Primary.ID)
		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("Error waiting for JIT policy (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
