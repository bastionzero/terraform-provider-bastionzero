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
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo", "bar"}, []string{string(verbtype.Shell)}),
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
					resource.TestCheckTypeSetElemAttr(resourceName, "verbs.*", string(verbtype.Shell)),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccTargetConnectPolicy_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy policies.TargetConnectPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Tunnel)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Policies.DeleteTargetConnectPolicy(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTargetConnectPolicy_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName1, []string{"foo"}, []string{string(verbtype.Shell)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name: &rName1,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigBasic(rName2, []string{"foo"}, []string{string(verbtype.Shell)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name: &rName2,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_TargetUsers(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Shell)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{{Username: "foo"}},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"bar"}, []string{string(verbtype.Shell)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:        &rName,
						TargetUsers: &[]policies.TargetUser{{Username: "bar"}},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemAttr(resourceName, "target_users.*", "bar"),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Verbs(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Shell)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:  &rName,
						Verbs: &[]policies.Verb{{Type: verbtype.Shell}},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemAttr(resourceName, "verbs.*", string(verbtype.Shell)),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update verbs
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName, []string{"foo"}, []string{string(verbtype.Tunnel)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:  &rName,
						Verbs: &[]policies.Verb{{Type: verbtype.Tunnel}},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemAttr(resourceName, "verbs.*", string(verbtype.Tunnel)),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
	desc1 := "desc1"
	desc2 := "desc2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigDescription(rName, []string{"foo"}, []string{string(verbtype.Shell)}, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:        &rName,
						Description: &desc1,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigDescription(rName, []string{"foo"}, []string{string(verbtype.Shell)}, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:        &rName,
						Description: &desc2,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Verify setting to empty string clears
			{
				Config: testAccTargetConnectPolicyConfigDescription(rName, []string{"foo"}, []string{string(verbtype.Shell)}, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
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
	acctest.FindTwoUsersOrSkip(t, ctx, subject1, subject2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1, *subject2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1, *subject2},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject2.ID, "type": string(subject2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicySubjects(ctx, []policies.Subject{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "0"),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Groups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
	group1 := new(policies.Group)
	group2 := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindTwoGroupsOrSkip(t, ctx, group1, group2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigGroups(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigGroups(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1, *group2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1, *group2},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group2.ID, "name": string(group2.Name)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccTargetConnectPolicyConfigGroups(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyGroups(ctx, []policies.Group{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Environments(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
	env1 := new(policies.Environment)
	env2 := new(policies.Environment)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindTwoEnvironmentsOrSkip(t, ctx, env1, env2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigEnvironments(rName, []string{"foo"}, []string{string(verbtype.Shell)}, []string{env1.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigEnvironments(rName, []string{"foo"}, []string{string(verbtype.Shell)}, []string{env1.ID, env2.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1, *env2},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env1.ID),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env2.ID),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccTargetConnectPolicyConfigEnvironments(rName, []string{"foo"}, []string{string(verbtype.Shell)}, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "environments.#", "0"),
				),
			},
		},
	})
}

func TestAccTargetConnectPolicy_Targets(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
	target1 := new(policies.Target)
	target2 := new(policies.Target)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindTwoBzeroTargetsOrSkip(t, ctx, target1, target2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConnectPolicyConfigTargets(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyTargets(ctx, []policies.Target{*target1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{*target1},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
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
				Config: testAccTargetConnectPolicyConfigTargets(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyTargets(ctx, []policies.Target{*target1, *target2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{*target1, *target2},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "targets.*", map[string]string{"id": target1.ID, "type": string(target1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "targets.*", map[string]string{"id": target2.ID, "type": string(target2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccTargetConnectPolicyConfigTargets(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicyTargets(ctx, []policies.Target{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:    &rName,
						Targets: &[]policies.Target{},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "targets.#", "0"),
				),
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
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs))
}

func testAccTargetConnectPolicyConfigDescription(rName string, targetUsers []string, verbs []string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  description = %[4]q
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), description)
}

func testAccTargetConnectPolicyConfigSubjects(rName string, targetUsers []string, verbs []string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  subjects = %[4]s
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), subjects.String())
}

func testAccTargetConnectPolicyConfigGroups(rName string, targetUsers []string, verbs []string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  groups = %[4]s
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), groups.String())
}

func testAccTargetConnectPolicyConfigEnvironments(rName string, targetUsers []string, verbs []string, environments []string) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  environments = %[4]s
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), acctest.ToTerraformStringList(environments))
}

func testAccTargetConnectPolicyConfigTargets(rName string, targetUsers []string, verbs []string, targets types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  targets = %[4]s
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), targets.String())
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
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.TargetConnect)),
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
