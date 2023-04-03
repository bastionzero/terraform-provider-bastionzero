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
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/organization"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
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

// SkipIfNotInAcceptanceTestMode performs the same check that the
// terraform-plugin-testing library performs to see if the test should be
// executed or not. Its logic is duplicated here, so we can call it ourselves
// before the Test() block in case there are additional things that need to be
// done that cannot be done in PreConfig() or PreCheck()
func SkipIfNotInAcceptanceTestMode(t *testing.T) {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc))
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

// ConfigCompose can be called to concatenate multiple strings to build test
// configurations
func ConfigCompose(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
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

// ListOrSetCount returns the number of elements in a list or set attribute
func ListOrSetCount(resourceState *terraform.ResourceState, listOrSetAttributeName string) (int, error) {
	rawCount, ok := resourceState.Primary.Attributes[fmt.Sprintf("%s.#", listOrSetAttributeName)]
	if !ok {
		return 0, fmt.Errorf("Could not find list/set attribute %s", listOrSetAttributeName)
	}
	count, err := strconv.Atoi(rawCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckListOrSetHasElements attempts to load a resource/datasource with name
// namedTFResource from the TF state, and then check that the list/set at
// listOrSetAttributeName has at least 1 element.
func CheckListOrSetHasElements(namedTFResource, listOrSetAttributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]

		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		total, err := ListOrSetCount(rs, listOrSetAttributeName)
		if err != nil {
			return err
		}

		if total < 1 {
			return fmt.Errorf("No %s retrieved", listOrSetAttributeName)
		}

		return nil
	}
}

// CheckTypeSetElemNestedAttrsFromResource ensures a subset map of values is
// stored in state for the given name (nameSecond) and key (attr) combination of
// attributes nested under a list or set block. The expected subset map is built
// by copying one for one the key and value combinations found at nameFirst in
// the state.
//
// Optionally, copy certain key and value cominations by providing a whitelist
// of keys. Otherwise, if keys list is empty, it is assumed all key and value
// pairs should be asserted to exist in one of the nested objects under a list
// or set block (specified by attr).
func CheckTypeSetElemNestedAttrsFromResource(nameFirst string, keys []string, nameSecond string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[nameFirst]
		if !ok {
			return fmt.Errorf("resource not found: %s", nameFirst)
		}

		var values map[string]string
		if len(keys) > 0 {
			// Build dictionary of keys to filter for in resource
			keysMap := make(map[string]struct{})
			for _, v := range keys {
				keysMap[v] = struct{}{}
			}

			// Create expected, nested object using only select keys from
			// provided resource
			values = make(map[string]string, 0)
			for k, v := range rs.Primary.Attributes {
				if _, ok := keysMap[k]; ok {
					values[k] = v
				}
			}
		} else {
			// Otherwise, assume all key and value pairs to exist in a nested
			// object
			values = rs.Primary.Attributes
		}

		return resource.TestCheckTypeSetElemNestedAttrs(nameSecond, attr, values)(s)
	}
}

// CheckResourceDisappears loads namedTFResource from the Terraform state and
// runs f to delete the API object at BastionZero. The ID passed to f is taken
// from the state file
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

// CheckAllPoliciesHaveSubjectID checks that all policies have at least one
// subject that matches an expected ID. It is expected that namedTFResource has
// a nested list/set attribute named "policies" and the container must contain
// another attribute named "subjects" that contains nested subject objects
func CheckAllPoliciesHaveSubjectID(namedTFResource, expectedSubjectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		totalPolicies, err := ListOrSetCount(rs, "policies")
		if err != nil {
			return err
		}

		if totalPolicies == 0 {
			return fmt.Errorf("list of policies is empty")
		}

		for i := 0; i < totalPolicies; i++ {
			if err := resource.TestCheckTypeSetElemNestedAttrs(namedTFResource, fmt.Sprintf("policies.%v.subjects.*", i), map[string]string{"id": expectedSubjectID})(s); err != nil {
				// This policy does not have at least one subject with a
				// matching ID.
				return err
			}
		}

		return nil
	}
}

// CheckAllPoliciesHaveGroupID checks that all policies have at least one group
// that matches an expected ID. It is expected that namedTFResource has a nested
// list/set attribute named "policies" and the container must contain another
// attribute named "groups" that contains nested group objects
func CheckAllPoliciesHaveGroupID(namedTFResource, expectedGroupID string) resource.TestCheckFunc {
	// TODO-Yuval: Potentially refactor and abstract common code with
	// CheckAllPoliciesHaveSubjectID()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[namedTFResource]
		if !ok {
			return fmt.Errorf("Not found: %s", namedTFResource)
		}

		totalPolicies, err := ListOrSetCount(rs, "policies")
		if err != nil {
			return err
		}

		if totalPolicies == 0 {
			return fmt.Errorf("list of policies is empty")
		}

		for i := 0; i < totalPolicies; i++ {
			if err := resource.TestCheckTypeSetElemNestedAttrs(namedTFResource, fmt.Sprintf("policies.%v.groups.*", i), map[string]string{"id": expectedGroupID})(s); err != nil {
				// This policy does not have at least one group with a matching
				// ID.
				return err
			}
		}

		return nil
	}
}

func typeof(v interface{}) string {
	// Source: https://stackoverflow.com/a/27160765
	return fmt.Sprintf("%T", v)
}

// FindNAPIObjectsOrSkip calls f to find a list of API objects at BastionZero
// and sets a variadic number (n) of pointers to the first n API objects found.
// The API object is converted to another type, MappedT, by calling mapF (pass
// the identity function if you don't want to map).
//
// If there are less than n API objects, then the current test is skipped.
func FindNAPIObjectsOrSkip[APIObject any, MappedT any](
	t *testing.T,
	f func(client *bzapi.Client, ctx context.Context) ([]APIObject, *http.Response, error),
	mapF func(apiObject APIObject) MappedT,
	mappedPointers ...*MappedT,
) {
	apiObjects, _, err := f(APIClient, context.Background())
	if err != nil {
		t.Fatalf("failed to list %v API objects: %s", typeof(new(APIObject)), err)
	}

	if len(apiObjects) < len(mappedPointers) {
		t.Skipf("skipping %s because we need at least %v %v API objects to test correctly but have %v", t.Name(), len(mappedPointers), typeof(new(APIObject)), len(apiObjects))
	}

	for i, mappedPointer := range mappedPointers {
		*mappedPointer = mapF(apiObjects[i])
	}
}

// FindNUsersOrSkip lists the users in the BastionZero organization and sets
// subjects to the first n subjects found. If there are less than n users, then
// the current test is skipped.
func FindNUsersOrSkip(t *testing.T, subjects ...*policies.Subject) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]users.User, *http.Response, error) {
		return client.Users.ListUsers(ctx)
	}, func(u users.User) policies.Subject {
		return policies.Subject{ID: u.ID, Type: u.GetSubjectType()}
	}, subjects...)
}

// FindNGroupsOrSkip lists the groups in the BastionZero organization and sets
// groups to the first n groups found. If there are less than n groups, then the
// current test is skipped.
func FindNGroupsOrSkip(t *testing.T, groups ...*policies.Group) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]organization.Group, *http.Response, error) {
		return client.Organization.ListGroups(ctx)
	}, func(g organization.Group) policies.Group {
		return policies.Group{ID: g.ID, Name: g.Name}
	}, groups...)
}

// FindNEnvironmentsOrSkip lists the environments in the BastionZero
// organization and sets envs to the first n environments found. If there are
// less than n environments, then the current test is skipped.
func FindNEnvironmentsOrSkip(t *testing.T, envs ...*policies.Environment) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]environments.Environment, *http.Response, error) {
		return client.Environments.ListEnvironments(ctx)
	}, func(e environments.Environment) policies.Environment {
		return policies.Environment{ID: e.ID}
	}, envs...)
}

// FindNBzeroTargetsOrSkip lists the Bzero targets in the BastionZero
// organization and sets bzeroTargets to the first n Bzero targets found. If
// there are less than n Bzero targets, then the current test is skipped.
func FindNBzeroTargetsOrSkip(t *testing.T, bzeroTargets ...*policies.Target) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.BzeroTarget, *http.Response, error) {
		return client.Targets.ListBzeroTargets(ctx)
	}, func(t targets.BzeroTarget) policies.Target {
		return policies.Target{ID: t.ID, Type: t.GetTargetType()}
	}, bzeroTargets...)
}

// FindNClusterTargetsOrSkip lists the Cluster targets in the BastionZero
// organization and sets clusterTargets to the first n Cluster targets found. If
// there are less than n Cluster targets, then the current test is skipped.
func FindNClusterTargetsOrSkip(t *testing.T, clusterTargets ...*policies.Cluster) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.ClusterTarget, *http.Response, error) {
		return client.Targets.ListClusterTargets(ctx)
	}, func(t targets.ClusterTarget) policies.Cluster {
		return policies.Cluster{ID: t.ID}
	}, clusterTargets...)
}

// FindNDbTargetsOrSkip lists the Db targets in the BastionZero organization and
// sets dbTargets to the first n Db targets found. If there are less than n Db
// targets, then the current test is skipped.
func FindNDbTargetsOrSkip(t *testing.T, dbTargets ...*policies.Target) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.DatabaseTarget, *http.Response, error) {
		return client.Targets.ListDatabaseTargets(ctx)
	}, func(t targets.DatabaseTarget) policies.Target {
		return policies.Target{ID: t.ID, Type: t.GetTargetType()}
	}, dbTargets...)
}

func ToTerraformStringList(arr []string) string {
	// Source: https://stackoverflow.com/questions/24489384/how-to-print-the-values-of-slices#comment126502244_53672500
	return strings.ReplaceAll(fmt.Sprintf("%+q", arr), "\" \"", "\",\"")
}
