package environment_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEnvironment_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_environment.test"
	var env environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccEnvironmentConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					// Check environment exists at BastionZero
					testAccCheckEnvironmentExists(resourceName, &env),
					// Check environment stored at BastionZero looks correct
					testAccCheckEnvironmentAttributes(&env, &expectedEnvironment{
						Name:                       &rName,
						Description:                bastionzero.PtrTo(""),
						OfflineCleanupTimeoutHours: bastionzero.PtrTo(environment.DefaultOfflineCleanupTimeoutHours)},
					),
					// Check computed values in TF state are correct
					testAccCheckResourceEnvironmentComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
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

func TestAccEnvironment_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_environment.test"
	var env environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Environments.DeleteEnvironment(ctx, id)
					}),
				),
				// The resource was deleted in CheckResourceDisappears (if no
				// error occurs in calling the API); therefore, the final plan
				// should not be empty (it should ask to re-create the object
				// since it was deleted by someone else).
				//
				// See:
				// https://developer.hashicorp.com/terraform/plugin/testing/testing-patterns#built-in-patterns
				// and
				// https://github.com/hashicorp/terraform-provider-aws/blob/main/docs/running-and-writing-acceptance-tests.md#disappears-acceptance-tests
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEnvironment_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	desc1 := "description1"
	desc2 := "description2"
	resourceName := "bastionzero_environment.test"
	var env1, env2 environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config that sets description
			{
				Config: testAccEnvironmentConfigDescription(rName, desc1),
				Check: resource.ComposeTestCheckFunc(
					// Check environment exists at BastionZero
					testAccCheckEnvironmentExists(resourceName, &env1),
					// Check environment stored at BastionZero looks correct
					testAccCheckEnvironmentAttributes(&env1, &expectedEnvironment{
						Name:                       &rName,
						Description:                bastionzero.PtrTo(desc1),
						OfflineCleanupTimeoutHours: bastionzero.PtrTo(environment.DefaultOfflineCleanupTimeoutHours)},
					),
					// Check computed values in TF state are correct
					testAccCheckResourceEnvironmentComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc1),
					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
				),
			},
			// Verify import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update description
			{
				Config: testAccEnvironmentConfigDescription(rName, desc2),
				Check: resource.ComposeTestCheckFunc(
					// Check environment exists at BastionZero
					testAccCheckEnvironmentExists(resourceName, &env2),
					// Check environment stored at BastionZero looks correct
					testAccCheckEnvironmentAttributes(&env2, &expectedEnvironment{
						Name:                       &rName,
						Description:                bastionzero.PtrTo(desc2),
						OfflineCleanupTimeoutHours: bastionzero.PtrTo(environment.DefaultOfflineCleanupTimeoutHours)},
					),
					testAccCheckResourceEnvironmentComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
				),
			},
		},
	})
}

func TestAccEnvironment_OfflineCleanupTimeoutHours(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	timeout1 := 24
	timeout2 := 48
	resourceName := "bastionzero_environment.test"
	var env1, env2 environments.Environment

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfigOfflineCleanupTimeoutHours(rName, strconv.Itoa(timeout1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env1),
					testAccCheckEnvironmentAttributes(&env1, &expectedEnvironment{
						Name:                       &rName,
						Description:                bastionzero.PtrTo(""),
						OfflineCleanupTimeoutHours: bastionzero.PtrTo(timeout1)},
					),
					testAccCheckResourceEnvironmentComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(timeout1)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEnvironmentConfigOfflineCleanupTimeoutHours(rName, strconv.Itoa(timeout2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &env2),
					testAccCheckEnvironmentAttributes(&env2, &expectedEnvironment{
						Name:                       &rName,
						Description:                bastionzero.PtrTo(""),
						OfflineCleanupTimeoutHours: bastionzero.PtrTo(timeout2)},
					),
					testAccCheckResourceEnvironmentComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(timeout2)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccEnvironment_RecreateOnNameChange(t *testing.T) {
	var afterCreate, afterUpdate environments.Environment
	name := acctest.RandomName()
	name2 := acctest.RandomName()
	resourceName := "bastionzero_environment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(context.Background(), t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfigName(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &afterCreate),
					testAccCheckEnvironmentAttributes(&afterCreate, &expectedEnvironment{Name: &name}),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			// Change name should force re-create the resource
			{
				Config: testAccEnvironmentConfigName(name2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists(resourceName, &afterUpdate),
					testAccCheckEnvironmentAttributes(&afterUpdate, &expectedEnvironment{Name: &name2}),
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					testAccCheckEnvironmentRecreated(t, &afterCreate, &afterUpdate),
				),
			},
		},
	})
}

func testAccEnvironmentConfigName(rName string) string {
	return fmt.Sprintf(`
resource "bastionzero_environment" "test" {
  name = %[1]q
}
`, rName)
}

func testAccEnvironmentConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_environment" "test" {
  description = %[2]q
  name = %[1]q
}
`, rName, description)
}

func testAccEnvironmentConfigOfflineCleanupTimeoutHours(rName string, timeoutHours string) string {
	return fmt.Sprintf(`
resource "bastionzero_environment" "test" {
  offline_cleanup_timeout_hours = %[2]q
  name = %[1]q
}
`, rName, timeoutHours)
}

type expectedEnvironment struct {
	Name                       *string
	Description                *string
	OfflineCleanupTimeoutHours *int
}

func testAccCheckEnvironmentAttributes(env *environments.Environment, expected *expectedEnvironment) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if expected.Name != nil && *expected.Name != env.Name {
			return fmt.Errorf("Bad name, expected \"%s\", got: %#v", *expected.Name, env.Name)
		}
		if expected.Description != nil && *expected.Description != env.Description {
			return fmt.Errorf("Bad description, expected \"%s\", got: %#v", *expected.Description, env.Description)
		}
		if expected.OfflineCleanupTimeoutHours != nil && *expected.OfflineCleanupTimeoutHours != int(env.OfflineCleanupTimeoutHours) {
			return fmt.Errorf("Bad offline_cleanup_timeout_hours, expected \"%v\", got: %#v", *expected.OfflineCleanupTimeoutHours, env.OfflineCleanupTimeoutHours)
		}

		return nil
	}
}

// testAccCheckEnvironmentExists checks that namedTFResource exists in the
// Terraform state and its ID represents an environment that exists at
// BastionZero. If the environment is found, its value is stored at the provided
// pointer.
func testAccCheckEnvironmentExists(namedTFResource string, environment *environments.Environment) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, environment, func(c *bastionzero.Client, ctx context.Context, id string) (*environments.Environment, *http.Response, error) {
		return c.Environments.GetEnvironment(ctx, id)
	})
}

// testAccCheckResourceEnvironmentComputedAttr checks all computed (read-only)
// attributes of an environment resource match expected values
func testAccCheckResourceEnvironmentComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
		resource.TestCheckResourceAttr(resourceName, "targets.%", "0"),
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestMatchResourceAttr(resourceName, "organization_id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestMatchResourceAttr(resourceName, "time_created", regexp.MustCompile(acctest.RFC3339RegexPattern)),
	)
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bastionzero_environment" {
			continue
		}

		// Try to find the environment
		_, _, err := acctest.APIClient.Environments.GetEnvironment(context.Background(), rs.Primary.ID)
		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("Error waiting for environment (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckEnvironmentRecreated(t *testing.T, before, after *environments.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID == after.ID {
			t.Fatalf("Expected change of environment IDs, but both were %v", before.ID)
		}
		return nil
	}
}
