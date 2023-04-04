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
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
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
			"child_policies.*",
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
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:                  &rName,
						Description:           bastionzero.PtrTo(""),
						Subjects:              &[]policies.Subject{},
						Groups:                &[]policies.Group{},
						ChildPolicies:         &asChildPolicies,
						AutomaticallyApproved: bastionzero.PtrTo(false),
						Duration:              bastionzero.PtrTo(uint(60)),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
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

func TestAccJITPolicy_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Policies.DeleteJITPolicy(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccJITPolicy_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigBasic(rName1, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name: &rName1,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update name
			{
				Config: testAccJITPolicyConfigBasic(rName2, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name: &rName2,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccJITPolicy_ChildPolicies(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)
	asChildPolicies := convertToChildPolicies(tcPolicy, kubePolicy, proxyPolicy)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:          &rName,
						ChildPolicies: bastionzero.PtrTo(asChildPolicies[:1]),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					testAccCheckResourceJITPolicyChildPolicies(resourceName, asChildPolicies[:1]...),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update child policies
			{
				Config: testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:          &rName,
						ChildPolicies: bastionzero.PtrTo(asChildPolicies[:2]),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					testAccCheckResourceJITPolicyChildPolicies(resourceName, asChildPolicies[:2]...),
				),
			},
			// Add another child policy
			{
				Config: testAccJITPolicyConfigBasic(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:          &rName,
						ChildPolicies: &asChildPolicies,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					testAccCheckResourceJITPolicyChildPolicies(resourceName, asChildPolicies...),
				),
			},
		},
	})
}

func TestAccJITPolicy_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)
	desc1 := "desc1"
	desc2 := "desc2"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigDescription(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:        &rName,
						Description: &desc1,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc1),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update description
			{
				Config: testAccJITPolicyConfigDescription(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:        &rName,
						Description: &desc2,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Verify setting to empty string clears
			{
				Config: testAccJITPolicyConfigDescription(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccJITPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var p policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)
	subject1 := new(policies.Subject)
	subject2 := new(policies.Subject)

	// Find two users or skip this entire test
	acctest.FindNUsersOrSkip(t, subject1, subject2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigSubjects(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update subjects
			{
				Config: testAccJITPolicyConfigSubjects(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1, *subject2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1, *subject2},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject2.ID, "type": string(subject2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccJITPolicyConfigSubjects(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicySubjects(ctx, []policies.Subject{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "0"),
				),
			},
		},
	})
}

func TestAccJITPolicy_Groups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var p policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)
	group1 := new(policies.Group)
	group2 := new(policies.Group)

	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group1, group2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigGroups(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update groups
			{
				Config: testAccJITPolicyConfigGroups(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1, *group2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1, *group2},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group2.ID, "name": string(group2.Name)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccJITPolicyConfigGroups(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, policy.FlattenPolicyGroups(ctx, []policies.Group{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &p),
					testAccCheckJITPolicyAttributes(t, &p, &expectedJITPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{},
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccJITPolicy_AutoApproved(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigAutoApproved(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:                  &rName,
						AutomaticallyApproved: bastionzero.PtrTo(true),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_approved", "true"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update auto_approved
			{
				Config: testAccJITPolicyConfigAutoApproved(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:                  &rName,
						AutomaticallyApproved: bastionzero.PtrTo(false),
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_approved", "false"),
				),
			},
		},
	})
}

func TestAccJITPolicy_Duration(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	duration1 := uint(20)
	duration2 := uint(40)
	resourceName := "bastionzero_jit_policy.test"
	var policy policies.JITPolicy
	tcPolicy, kubePolicy, proxyPolicy := getChildPoliciesOrSkip(ctx, t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckJITPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJITPolicyConfigDuration(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, duration1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:     &rName,
						Duration: &duration1,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "duration", strconv.Itoa(int(duration1))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJITPolicyConfigDuration(rName, []string{tcPolicy.ID, kubePolicy.ID, proxyPolicy.ID}, duration2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJITPolicyExists(resourceName, &policy),
					testAccCheckJITPolicyAttributes(t, &policy, &expectedJITPolicy{
						Name:     &rName,
						Duration: &duration2,
					}),
					testAccCheckResourceJITPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "duration", strconv.Itoa(int(duration2))),
				),
			},
		},
	})
}

func TestJITPolicy_InvalidChildPolicies(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty child policies not permitted
				Config:      testAccJITPolicyConfigBasic("test", []string{}),
				ExpectError: regexp.MustCompile(`at least 1 elements`),
			},
		},
	})
}

func TestJITPolicy_InvalidDuration(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid duration not permitted
				Config:      testAccJITPolicyConfigDuration("test", []string{"foo"}, 0),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestJITPolicy_InvalidName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty name not permitted
				Config:      testAccJITPolicyConfigBasic("", []string{"foo"}),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestJITPolicy_InvalidSubjects(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid subject type not permitted
				Config:      testAccJITPolicyConfigSubjects("test", []string{"foo"}, policy.FlattenPolicySubjects(context.Background(), []policies.Subject{{ID: "foo", Type: "foo"}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccJITPolicyConfigBasic(rName string, childPolicyIDs []string) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String())
}

func testAccJITPolicyConfigDescription(rName string, childPolicyIDs []string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  description = %[3]q
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), description)
}

func testAccJITPolicyConfigSubjects(rName string, childPolicyIDs []string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_jit_policy" "test" {
  subjects = %[3]s
  name = %[1]q
  child_policies = %[2]s
}
`, rName, toChildPoliciesSet(childPolicyIDs).String(), subjects.String())
}

func testAccJITPolicyConfigGroups(rName string, childPolicyIDs []string, groups types.Set) string {
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

type expectedJITPolicy struct {
	Name                  *string
	Description           *string
	Subjects              *[]policies.Subject
	Groups                *[]policies.Group
	ChildPolicies         *[]policies.ChildPolicy
	AutomaticallyApproved *bool
	Duration              *uint
}

func testAccCheckJITPolicyAttributes(t *testing.T, policy *policies.JITPolicy, expected *expectedJITPolicy) resource.TestCheckFunc {
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

func testAccCheckJITPolicyExists(namedTFResource string, policy *policies.JITPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.JITPolicy, *http.Response, error) {
		return c.Policies.GetJITPolicy(ctx, id)
	})
}

func testAccCheckResourceJITPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.JustInTime)),
	)
}

func testAccCheckJITPolicyDestroy(s *terraform.State) error {
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
