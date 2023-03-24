package bztftest

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for
	// example assertions about the appropriate environment variables being set
	// are common to see in a pre-check function.

	if apiSecret := os.Getenv("BASTIONZERO_API_SECRET"); apiSecret == "" {
		t.Fatal("The BASTIONZERO_API_SECRET environment variable must be set for acceptance tests.")
	}
}

// TestNamePrefix is a prefix for randomly generated names used during
// acceptance testing
const TestNamePrefix = "tf-acc-test-"

func RandomTestName(additionalNames ...string) string {
	prefix := TestNamePrefix
	for _, n := range additionalNames {
		prefix += "-" + strings.Replace(n, " ", "_", -1)
	}
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(10))
}

func TestAccCheckEnvironmentExists(n string, environment *environments.Environment) resource.TestCheckFunc {
	return TestAccCheckExistsAtBastionZero(n, environment, func(c *bastionzero.Client, ctx context.Context, s string) (*environments.Environment, error) {
		foundEnv, _, err := c.Environments.GetEnvironment(ctx, s)
		return foundEnv, err
	})
}

// TestAccCheckExistsAtBastionZero attempts to load a resource/datasource with
// name namedTFResource from the TF state and find an API object at BastionZero,
// using f, with the resource's ID.
//
// The provided pointer is set if there is no error when calling BastionZero. It
// can be examined to check that what exists at BastionZero matches what is
// actually set in the TF config/state.
func TestAccCheckExistsAtBastionZero[T any](namedTFResource string, apiObject *T, f func(*bastionzero.Client, context.Context, string) (*T, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Load from state
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set in the loaded resource")
		}

		client := GetBastionZeroClient()

		// Try to find the API object
		foundApiObject, err := f(client, context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*apiObject = *foundApiObject

		return nil
	}
}
