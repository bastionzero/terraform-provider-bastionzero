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
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
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
				Config: testAccTargetConnectPolicyConfigDescription(rName, []string{"foo"}, []string{string(verbtype.Tunnel)}, desc2),
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
		},
	})
}

func TestAccTargetConnectPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_targetconnect_policy.test"
	var policy1, policy2 policies.TargetConnectPolicy
	subjects1 := new([]policies.Subject)
	subjects2 := new([]policies.Subject)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { findTwoUsersOrSkip(t, ctx, subjects1, subjects2) },
				Config:    testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Shell)}, policy.FlattenPolicySubjects(ctx, *subjects1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy1),
					testAccCheckTargetConnectPolicyAttributes(t, &policy1, &expectedTargetConnectPolicy{
						Name:     &rName,
						Subjects: subjects1,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": "id", "type": string(subjecttype.User)}),
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
				Config: testAccTargetConnectPolicyConfigSubjects(rName, []string{"foo"}, []string{string(verbtype.Tunnel)}, policy.FlattenPolicySubjects(ctx, *subjects2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy2),
					testAccCheckTargetConnectPolicyAttributes(t, &policy2, &expectedTargetConnectPolicy{
						Name:     &rName,
						Subjects: subjects2,
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": "id2", "type": string(subjecttype.User)}),
				),
			},
		},
	})
}

func findTwoUsersOrSkip(t *testing.T, ctx context.Context, subjects1, subjects2 *[]policies.Subject) {
	client := acctest.APIClient

	users, _, err := client.Users.ListUsers(ctx)
	if err != nil {
		t.Fatalf("failed to list users: %s", err)
	}

	if len(users) < 200 {
		t.Skipf("skipping %s because we need at least two users to test correctly but have %v", t.Name(), len(users))
	}

	*subjects1 = []policies.Subject{{ID: users[0].ID, Type: users[0].GetSubjectType()}}
	*subjects2 = []policies.Subject{{ID: users[1].ID, Type: users[1].GetSubjectType()}}
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
	// 	 targets = [
	//    {
	//      id   = "cb0aecd0-2aae-4b2b-acda-5197250f1851",
	//      type = "Bzero"
	//    }
	//  ]

	return fmt.Sprintf(`
resource "bastionzero_targetconnect_policy" "test" {
  subjects = %[4]s
  name = %[1]q
  target_users = %[2]s
  verbs = %[3]s
}
`, rName, acctest.ToTerraformStringList(targetUsers), acctest.ToTerraformStringList(verbs), subjects.String())
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
