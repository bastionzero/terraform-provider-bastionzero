package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	bzapi "github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const RFC3339RegexPattern = `^[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9]+)?([Zz]|([+-]([01][0-9]|2[0-3]):[0-5][0-9]))$`
const UUIDV4RegexPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

var (
	// testAccAPIClientConfigure ensures APIClient is only configured once
	testAccAPIClientConfigure sync.Once

	// APIClient is a BastionZero API client.
	//
	// This can be used in testing code for API calls without requiring the use
	// of saving and referencing specific ProviderFactories instances.
	//
	// PreCheck(t) must be called before using this.
	APIClient *bzapi.Client

	// TestProtoV6ProviderFactories are used to instantiate a provider during
	// testing. The factory function will be invoked for every Terraform CLI
	// command executed to create a provider server to which the CLI can
	// reattach.
	TestProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
)

func init() {
	TestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"bastionzero": providerserver.NewProtocol6WithError(bastionzero.New("test")()),
	}

}

// PreCheck verifies and sets required provider testing configuration
//
// PreCheck makes assertions that must hold true in order to run an acceptance
// test. The test fails immediately if any of these assertions fails.
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration and ensures testing
// functions that attempt to call BastionZero APIs directly via APIClient are
// previously configured.
func PreCheck(ctx context.Context, t *testing.T) {
	testAccAPIClientConfigure.Do(func() {
		// You can add code here to run prior to any test case execution, for
		// example assertions about the appropriate environment variables being
		// set are common to see in a pre-check function.
		if apiSecret := os.Getenv("BASTIONZERO_API_SECRET"); apiSecret == "" {
			t.Fatal("The BASTIONZERO_API_SECRET environment variable must be set in order to run acceptance tests.")
		}

		// Create dummy provider so that we can access a properly configured
		// BastionZero client and test provider configuration e2e
		testProvider := bastionzero.New("test")()

		// Get schema from the provider
		schemaResponse := new(provider.SchemaResponse)
		testProvider.Schema(ctx, provider.SchemaRequest{}, schemaResponse)

		// Create empty config
		testConfig := tfsdk.Config{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"api_endpoint": tftypes.String,
					"api_secret":   tftypes.String,
				},
			}, map[string]tftypes.Value{
				"api_endpoint": tftypes.NewValue(tftypes.String, nil),
				"api_secret":   tftypes.NewValue(tftypes.String, nil),
			}),
			Schema: schemaResponse.Schema,
		}

		// Call Configure on the provider
		configureResponse := new(provider.ConfigureResponse)
		terraform.NewResourceConfigRaw(nil)
		testProvider.Configure(ctx, provider.ConfigureRequest{Config: testConfig}, configureResponse)

		// Parse the API client and save
		apiClient, ok := configureResponse.ResourceData.(*bzapi.Client)
		if !ok {
			t.Fatalf("expected provider to contain a *bastionzero.Client in its ResourceData")
		}
		APIClient = apiClient
	})
}

// TestNamePrefix is a prefix for randomly generated names used during
// acceptance testing
const TestNamePrefix = "tf-acc-test-"

// RandomName creates a random name suitable for named BastionZero API objects
// that are created during acceptance tests.
func RandomName(additionalNames ...string) string {
	prefix := TestNamePrefix
	for _, n := range additionalNames {
		prefix += "-" + strings.Replace(n, " ", "_", -1)
	}
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(10))
}

// CheckExistsAtBastionZero attempts to load a resource/datasource with name
// namedTFResource from the TF state and find an API object at BastionZero,
// using f, with the resource's ID.
//
// The provided pointer is set if there is no error when calling BastionZero. It
// can be examined to check that what exists at BastionZero matches what is
// actually set in the TF config/state.
func CheckExistsAtBastionZero[T any](namedTFResource string, apiObject *T, f func(client *bzapi.Client, ctx context.Context, id string) (*T, *http.Response, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", namedTFResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", namedTFResource)
		}

		// Try to find the API object
		foundApiObject, _, err := f(APIClient, context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*apiObject = *foundApiObject

		return nil
	}
}

// CheckListHasElements attempts to load a resource/datasource with name
// namedTFResource from the TF state, and then check that the list at
// listAttributeName has at least 1 element.
func CheckListHasElements(namedTFResource, listAttributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]

		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		rawTotal, ok := rs.Primary.Attributes[fmt.Sprintf("%s.#", listAttributeName)]
		if !ok {
			return fmt.Errorf("Not found %s", listAttributeName)
		}

		total, err := strconv.Atoi(rawTotal)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No %s retrieved", listAttributeName)
		}

		return nil
	}
}

func CheckResourceDisappears(namedTFResource string, f func(client *bzapi.Client, ctx context.Context, id string) (*http.Response, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", namedTFResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", namedTFResource)
		}

		// Try to delete the API object
		_, err := f(APIClient, context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}
