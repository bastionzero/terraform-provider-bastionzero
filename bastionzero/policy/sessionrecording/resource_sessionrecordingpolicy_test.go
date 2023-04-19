package sessionrecording_test

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
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccSessionRecordingPolicy_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccSessionRecordingPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
						Subjects:    &[]policies.Subject{},
						Groups:      &[]policies.Group{},
						RecordInput: bastionzero.PtrTo(false),
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "record_input", "false"),
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

func TestAccSessionRecordingPolicy_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Policies.DeleteSessionRecordingPolicy(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSessionRecordingPolicy_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigBasic(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name: &rName1,
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
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
				Config: testAccSessionRecordingPolicyConfigBasic(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name: &rName2,
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccSessionRecordingPolicy_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy
	desc1 := "desc1"
	desc2 := "desc2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigDescription(rName, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						Description: &desc1,
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
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
				Config: testAccSessionRecordingPolicyConfigDescription(rName, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						Description: &desc2,
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Verify setting to empty string clears
			{
				Config: testAccSessionRecordingPolicyConfigDescription(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccSessionRecordingPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var p policies.SessionRecordingPolicy
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
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
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
				Config: testAccSessionRecordingPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1, *subject2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1, *subject2},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject2.ID, "type": string(subject2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccSessionRecordingPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "0"),
				),
			},
		},
	})
}

func TestAccSessionRecordingPolicy_Groups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var p policies.SessionRecordingPolicy
	group1 := new(policies.Group)
	group2 := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkipAsPolicyGroup(t, group1, group2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
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
				Config: testAccSessionRecordingPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1, *group2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1, *group2},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group2.ID, "name": string(group2.Name)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccSessionRecordingPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &p),
					testAccCheckSessionRecordingPolicyAttributes(t, &p, &expectedSessionRecordingPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{},
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccSessionRecordingPolicy_RecordInput(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_sessionrecording_policy.test"
	var policy policies.SessionRecordingPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSessionRecordingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionRecordingPolicyConfigRecordInput(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						RecordInput: bastionzero.PtrTo(true),
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record_input", "true"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update record_input
			{
				Config: testAccSessionRecordingPolicyConfigRecordInput(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSessionRecordingPolicyExists(resourceName, &policy),
					testAccCheckSessionRecordingPolicyAttributes(t, &policy, &expectedSessionRecordingPolicy{
						Name:        &rName,
						RecordInput: bastionzero.PtrTo(false),
					}),
					testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record_input", "false"),
				),
			},
		},
	})
}

func TestSessionRecordingPolicy_InvalidName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty name not permitted
				Config:      testAccSessionRecordingPolicyConfigBasic(""),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestSessionRecordingPolicy_InvalidSubjects(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid subject type not permitted
				Config:      testAccSessionRecordingPolicyConfigSubjects("test", policy.FlattenPolicySubjects(context.Background(), []policies.Subject{{ID: uuid.New().String(), Type: "foo"}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				// Invalid ID not permitted
				Config:      testAccSessionRecordingPolicyConfigSubjects("test", policy.FlattenPolicySubjects(context.Background(), []policies.Subject{{ID: "foo", Type: subjecttype.User}})),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccSessionRecordingPolicyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  name = %[1]q
}
`, rName)
}

func testAccSessionRecordingPolicyConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  description = %[2]q
  name = %[1]q
}
`, rName, description)
}

func testAccSessionRecordingPolicyConfigSubjects(rName string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  subjects = %[2]s
  name = %[1]q
}
`, rName, subjects.String())
}

func testAccSessionRecordingPolicyConfigGroups(rName string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  groups = %[2]s
  name = %[1]q
}
`, rName, groups.String())
}

func testAccSessionRecordingPolicyConfigRecordInput(rName string, recordInput bool) string {
	return fmt.Sprintf(`
resource "bastionzero_sessionrecording_policy" "test" {
  record_input = %[2]t
  name = %[1]q
}
`, rName, recordInput)
}

type expectedSessionRecordingPolicy struct {
	Name        *string
	Description *string
	Subjects    *[]policies.Subject
	Groups      *[]policies.Group
	RecordInput *bool
}

func testAccCheckSessionRecordingPolicyAttributes(t *testing.T, policy *policies.SessionRecordingPolicy, expected *expectedSessionRecordingPolicy) resource.TestCheckFunc {
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
		if expected.RecordInput != nil && *expected.RecordInput != policy.GetRecordInput() {
			return fmt.Errorf("Bad record_input, expected \"%t\", got: %#v", *expected.RecordInput, policy.RecordInput)
		}

		return nil
	}
}

func testAccCheckSessionRecordingPolicyExists(namedTFResource string, policy *policies.SessionRecordingPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.SessionRecordingPolicy, *http.Response, error) {
		return c.Policies.GetSessionRecordingPolicy(ctx, id)
	})
}

func testAccCheckResourceSessionRecordingPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.SessionRecording)),
	)
}

func testAccCheckSessionRecordingPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bastionzero_sessionrecording_policy" {
			continue
		}

		// Try to find the policy
		_, _, err := acctest.APIClient.Policies.GetSessionRecordingPolicy(context.Background(), rs.Primary.ID)
		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("Error waiting for session recording policy (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
