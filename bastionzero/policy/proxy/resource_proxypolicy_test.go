package proxy_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccProxyPolicy_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccProxyPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:         &rName,
						Description:  bastionzero.PtrTo(""),
						Subjects:     &[]policies.Subject{},
						Groups:       &[]policies.Group{},
						Environments: &[]policies.Environment{},
						Targets:      &[]policies.Target{},
						TargetUsers:  &[]policies.TargetUser{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					// Check that unspecified values remain null
					resource.TestCheckNoResourceAttr(resourceName, "subjects"),
					resource.TestCheckNoResourceAttr(resourceName, "groups"),
					resource.TestCheckNoResourceAttr(resourceName, "environments"),
					resource.TestCheckNoResourceAttr(resourceName, "targets"),
					resource.TestCheckNoResourceAttr(resourceName, "target_users"),
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

func TestAccProxyPolicy_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Policies.DeleteProxyPolicy(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProxyPolicy_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigBasic(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name: &rName1,
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
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
				Config: testAccProxyPolicyConfigBasic(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name: &rName2,
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccProxyPolicy_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy
	desc1 := "desc1"
	desc2 := "desc2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigDescription(rName, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						Description: &desc1,
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
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
				Config: testAccProxyPolicyConfigDescription(rName, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						Description: &desc2,
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Verify setting to empty string clears
			{
				Config: testAccProxyPolicyConfigDescription(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccProxyPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var p policies.ProxyPolicy
	subject1 := new(policies.Subject)
	subject2 := new(policies.Subject)

	// These checks are here, instead of being inlined in PreCheck field,
	// because we need subject1 and subject2 to have values before using them as
	// arguments in the Test block below. Otherwise, any immediate pointer
	// dereference (e.g. in the TestSteps) will have the values set to nil which
	// is not what we want.
	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	// Find two users or skip this entire test
	acctest.FindNUsersOrSkipAsPolicySubject(t, subject1, subject2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
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
				Config: testAccProxyPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1, *subject2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1, *subject2},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject2.ID, "type": string(subject2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccProxyPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "0"),
				),
			},
		},
	})
}

func TestAccProxyPolicy_Groups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var p policies.ProxyPolicy
	group1 := new(policies.Group)
	group2 := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group1, group2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
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
				Config: testAccProxyPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1, *group2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1, *group2},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group2.ID, "name": string(group2.Name)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccProxyPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccProxyPolicy_Environments(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var p policies.ProxyPolicy
	env1 := new(policies.Environment)
	env2 := new(policies.Environment)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNEnvironmentsOrSkipAsPolicyEnvironment(t, env1, env2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigEnvironments(rName, []string{env1.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env1.ID),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update environments
			{
				Config: testAccProxyPolicyConfigEnvironments(rName, []string{env1.ID, env2.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1, *env2},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env1.ID),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env2.ID),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccProxyPolicyConfigEnvironments(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "environments.#", "0"),
				),
			},
		},
	})
}

func TestAccProxyPolicy_Targets(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var p policies.ProxyPolicy
	target1 := new(policies.Target)
	target2 := new(policies.Target)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNDbTargetsOrSkipAsPolicyTarget(t, target1, target2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigTargets(rName, policy.FlattenPolicyTargets(ctx, []policies.Target{*target1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{*target1},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "targets.*", map[string]string{"id": target1.ID, "type": string(target1.Type)}),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update targets
			{
				Config: testAccProxyPolicyConfigTargets(rName, policy.FlattenPolicyTargets(ctx, []policies.Target{*target1, *target2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{*target1, *target2},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "targets.*", map[string]string{"id": target1.ID, "type": string(target1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "targets.*", map[string]string{"id": target2.ID, "type": string(target2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccProxyPolicyConfigTargets(rName, policy.FlattenPolicyTargets(ctx, []policies.Target{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &p),
					testAccCheckProxyPolicyAttributes(t, &p, &expectedProxyPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "targets.#", "0"),
				),
			},
		},
	})
}

func TestAccProxyPolicy_TargetUsers(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_proxy_policy.test"
	var policy policies.ProxyPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProxyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxyPolicyConfigTargetUsers(rName, []string{"foo"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{{Username: "foo"}},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "target_users.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "foo"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update target users
			{
				Config: testAccProxyPolicyConfigTargetUsers(rName, []string{"bar"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{{Username: "bar"}},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "target_users.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "bar"),
				),
			},
			// Add another target user
			{
				Config: testAccProxyPolicyConfigTargetUsers(rName, []string{"bar", "baz"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{{Username: "bar"}, {Username: "baz"}},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "target_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "bar"),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "baz"),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccProxyPolicyConfigTargetUsers(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProxyPolicyExists(resourceName, &policy),
					testAccCheckProxyPolicyAttributes(t, &policy, &expectedProxyPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{},
					}),
					testAccCheckResourceProxyPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "target_users.#", "0"),
				),
			},
		},
	})
}

func TestAccProxyPolicy_MutualExclTargetsEnvs(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Cannot specify both environments and targets
				Config: `
				resource "bastionzero_proxy_policy" "test" {
			      name = "foo"
				  environments = []
				  targets = []
				}
				`,
				ExpectError: regexp.MustCompile(`cannot be configured together`),
			},
		},
	})
}

func TestProxyPolicy_InvalidTargets(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid target type not permitted
				Config:      testAccProxyPolicyConfigTargets("test", policy.FlattenPolicyTargets(context.Background(), []policies.Target{{ID: uuid.New().String(), Type: "foo"}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Invalid ID not permitted
				Config:      testAccProxyPolicyConfigTargets("test", policy.FlattenPolicyTargets(context.Background(), []policies.Target{{ID: "foo", Type: targettype.Bzero}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestProxyPolicy_InvalidName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty name not permitted
				Config:      testAccProxyPolicyConfigBasic(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestProxyPolicy_InvalidSubjects(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid subject type not permitted
				Config:      testAccProxyPolicyConfigSubjects("test", policy.FlattenPolicySubjects(context.Background(), []policies.Subject{{ID: uuid.New().String(), Type: "foo"}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Invalid ID not permitted
				Config:      testAccProxyPolicyConfigSubjects("test", policy.FlattenPolicySubjects(context.Background(), []policies.Subject{{ID: "foo", Type: subjecttype.User}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccProxyPolicyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  name = %[1]q
}
`, rName)
}

func testAccProxyPolicyConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  description = %[2]q
  name = %[1]q
}
`, rName, description)
}

func testAccProxyPolicyConfigSubjects(rName string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  subjects = %[2]s
  name = %[1]q
}
`, rName, subjects.String())
}

func testAccProxyPolicyConfigGroups(rName string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  groups = %[2]s
  name = %[1]q
}
`, rName, groups.String())
}

func testAccProxyPolicyConfigEnvironments(rName string, environments []string) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  environments = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(environments))
}

func testAccProxyPolicyConfigTargets(rName string, targets types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  targets = %[2]s
  name = %[1]q
}
`, rName, targets.String())
}

func testAccProxyPolicyConfigTargetUsers(rName string, targetUsers []string) string {
	return fmt.Sprintf(`
resource "bastionzero_proxy_policy" "test" {
  target_users = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(targetUsers))
}

type expectedProxyPolicy struct {
	Name         *string
	Description  *string
	Subjects     *[]policies.Subject
	Groups       *[]policies.Group
	Environments *[]policies.Environment
	Targets      *[]policies.Target
	TargetUsers  *[]policies.TargetUser
}

func testAccCheckProxyPolicyAttributes(t *testing.T, policy *policies.ProxyPolicy, expected *expectedProxyPolicy) resource.TestCheckFunc {
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

		return nil
	}
}

func testAccCheckProxyPolicyExists(namedTFResource string, policy *policies.ProxyPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.ProxyPolicy, *http.Response, error) {
		return c.Policies.GetProxyPolicy(ctx, id)
	})
}

func testAccCheckResourceProxyPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.Proxy)),
	)
}

func testAccCheckProxyPolicyDestroy(s *terraform.State) error {
	return acctest.CheckAllResourcesWithTypeDestroyed(
		"bastionzero_proxy_policy",
		func(client *bastionzero.Client, ctx context.Context, id string) (*policies.ProxyPolicy, *http.Response, error) {
			return client.Policies.GetProxyPolicy(ctx, id)
		},
	)(s)
}
