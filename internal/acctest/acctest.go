package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	bzapi "github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/organization"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"golang.org/x/exp/slices"
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

// FindBastionZeroAPIObjectFunc is a function that finds and returns a
// BastionZero API object given its ID
type FindBastionZeroAPIObjectFunc[T any] func(client *bzapi.Client, ctx context.Context, id string) (*T, *http.Response, error)

// CheckAllResourcesWithTypeDestroyed loops through all resources in the
// Terraform state and verifies that each resource of type resourceType is
// destroyed (no longer exists at BastionZero). Function f is used to find the
// API object given the resource's ID (parsed from Terraform state).
func CheckAllResourcesWithTypeDestroyed[T any](resourceType string, f FindBastionZeroAPIObjectFunc[T]) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			// Try to find the API object at BastionZero
			_, _, err := f(APIClient, context.Background(), rs.Primary.ID)

			// Check if we found the object
			if err == nil {
				return fmt.Errorf("Resource %s (id=%s) still exists at BastionZero", resourceType, rs.Primary.ID)
			}

			// We expect a 404 Not Found. If we receive any other error,
			// something might be wrong so return the error
			if !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
				return err
			}
		}

		// All resources of specified resourceType are gone
		return nil
	}
}

// CheckExistsAtBastionZero attempts to load a resource/datasource with name
// namedTFResource from the TF state and find an API object at BastionZero,
// using f, with the resource's ID.
//
// The provided pointer is set if there is no error when calling BastionZero. It
// can be examined to check that what exists at BastionZero matches what is
// actually set in the TF config/state.
func CheckExistsAtBastionZero[T any](namedTFResource string, apiObject *T, f FindBastionZeroAPIObjectFunc[T]) resource.TestCheckFunc {
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
// Optionally, copy certain key and value combinations by providing a whitelist
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

// CheckResourceAttrIsOneOf ensures the Terraform state value is equal to one of
// the supplied values.
func CheckResourceAttrIsOneOf(name, key string, values []string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key, func(value string) error {
		if slices.Contains(values, value) {
			return nil
		}
		return fmt.Errorf("%s: Attribute '%s' expected to be one of %#v, got %#v", name, key, values, value)
	})
}

func typeof(v interface{}) string {
	// Source: https://stackoverflow.com/a/27160765
	return fmt.Sprintf("%T", v)
}

func identity[T any](a T) T {
	return a
}

// FindNAPIObjectsOrSkip calls f to find a list of API objects at BastionZero
// and sets a variadic number (n) of pointers to the first n API objects found.
// The API object is converted to another type, MappedT, by calling mapF (pass
// the identity function if you don't want to map).
//
// Additionally, pass filterF if you wish to filter certain API objects from
// being included as candidates (pass nil if you don't want to filter).
//
// If there are less than n API objects, then the current test is skipped.
func FindNAPIObjectsOrSkip[APIObject any, MappedT any](
	t *testing.T,
	f func(client *bzapi.Client, ctx context.Context) ([]APIObject, *http.Response, error),
	mapF func(apiObject APIObject) MappedT,
	filterF func(apiObject APIObject) bool,
	mappedPointers ...*MappedT,
) {
	if mapF == nil {
		panic("mapF cannot be nil. Use the identity function if you don't want to map")
	}

	apiObjects, _, err := f(APIClient, context.Background())
	if err != nil {
		t.Fatalf("failed to list %v API objects: %s", typeof(new(APIObject)), err)
	}

	// Apply optional filter
	if filterF != nil {
		var filteredAPIObjects []APIObject
		for _, apiObject := range apiObjects {
			if filterF(apiObject) {
				filteredAPIObjects = append(filteredAPIObjects, apiObject)
			}
		}

		// Use filtered list
		apiObjects = filteredAPIObjects
	}

	if len(apiObjects) < len(mappedPointers) {
		t.Skipf("skipping %s because we need at least %v %v API objects to test correctly but have %v", t.Name(), len(mappedPointers), typeof(new(APIObject)), len(apiObjects))
	}

	for i, mappedPointer := range mappedPointers {
		*mappedPointer = mapF(apiObjects[i])
	}
}

// FindNUsersOrSkip lists the users in the BastionZero organization and sets
// userPointers to the first n users found. If there are less than n users, then
// the current test is skipped.
//
// If you need the users mapped as the policy type (policies.Subject), use
// FindNUsersOrSkipAsPolicySubject() instead.
func FindNUsersOrSkip(t *testing.T, userPointers ...*users.User) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]users.User, *http.Response, error) {
		return client.Users.ListUsers(ctx)
	}, identity[users.User], nil, userPointers...)
}

// FindNUsersOrSkipAsPolicySubject lists the users in the BastionZero
// organization and sets subjects to the first n subjects found. If there are
// less than n users, then the current test is skipped.
func FindNUsersOrSkipAsPolicySubject(t *testing.T, subjects ...*policies.Subject) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]users.User, *http.Response, error) {
		return client.Users.ListUsers(ctx)
	}, func(u users.User) policies.Subject {
		return policies.Subject{ID: u.ID, Type: u.GetSubjectType()}
	}, nil, subjects...)
}

// FindNServiceAccountsOrSkip lists the service accounts in the BastionZero
// organization and sets serviceAccounts to the first n service accounts found.
// If there are less than n service accounts, then the current test is skipped.
func FindNServiceAccountsOrSkip(t *testing.T, serviceAccounts ...*serviceaccounts.ServiceAccount) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]serviceaccounts.ServiceAccount, *http.Response, error) {
		return client.ServiceAccounts.ListServiceAccounts(ctx)
	}, identity[serviceaccounts.ServiceAccount], nil, serviceAccounts...)
}

// FindNGroupsOrSkip lists the groups in the BastionZero organization and sets
// groups to the first n groups found. If there are less than n groups, then the
// current test is skipped.
//
// If you need the groups mapped as the policy type (policies.Group), use
// FindNGroupsOrSkipAsPolicyGroup() instead.
func FindNGroupsOrSkip(t *testing.T, groups ...*organization.Group) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]organization.Group, *http.Response, error) {
		return client.Organization.ListGroups(ctx)
	}, identity[organization.Group], nil, groups...)
}

// FindNGroupsOrSkipAsPolicyGroup lists the groups in the BastionZero
// organization and sets groups to the first n groups found. If there are less
// than n groups, then the current test is skipped.
func FindNGroupsOrSkipAsPolicyGroup(t *testing.T, groups ...*policies.Group) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]organization.Group, *http.Response, error) {
		return client.Organization.ListGroups(ctx)
	}, func(g organization.Group) policies.Group {
		return policies.Group{ID: g.ID, Name: g.Name}
	}, nil, groups...)
}

// FindNEnvironmentsOrSkip lists the environments in the BastionZero
// organization and sets envs to the first n environments found. If there are
// less than n environments, then the current test is skipped.
//
// If you need the environments mapped as the policy type
// (policies.Environment), use FindNEnvironmentsOrSkipAsPolicyEnvironment()
// instead.
func FindNEnvironmentsOrSkip(t *testing.T, envs ...*environments.Environment) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]environments.Environment, *http.Response, error) {
		return client.Environments.ListEnvironments(ctx)
	}, identity[environments.Environment], func(e environments.Environment) bool {
		// IMPORTANT: We must filter out environments that are concurrently
		// being created by other parallel acceptance tests because they could
		// be deleted by the time the caller of this function uses them
		return !strings.HasPrefix(e.Name, TestNamePrefix)
	}, envs...)
}

// FindNEnvironmentsOrSkipAsPolicyEnvironment lists the environments in the
// BastionZero organization and sets envs to the first n environments found. If
// there are less than n environments, then the current test is skipped.
func FindNEnvironmentsOrSkipAsPolicyEnvironment(t *testing.T, envs ...*policies.Environment) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]environments.Environment, *http.Response, error) {
		return client.Environments.ListEnvironments(ctx)
	}, func(e environments.Environment) policies.Environment {
		return policies.Environment{ID: e.ID}
	}, func(e environments.Environment) bool {
		// IMPORTANT: We must filter out environments that are concurrently
		// being created by other parallel acceptance tests because they could
		// be deleted by the time the caller of this function uses them
		return !strings.HasPrefix(e.Name, TestNamePrefix)
	}, envs...)
}

// FindNBzeroTargetsOrSkip lists the Bzero targets in the BastionZero
// organization and sets bzeroTargets to the first n Bzero targets found. If
// there are less than n Bzero targets, then the current test is skipped.
//
// If you need the targets mapped as the policy type (policies.Target), use
// FindNBzeroTargetsOrSkipAsPolicyTarget() instead.
func FindNBzeroTargetsOrSkip(t *testing.T, bzeroTargets ...*targets.BzeroTarget) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.BzeroTarget, *http.Response, error) {
		return client.Targets.ListBzeroTargets(ctx)
	}, identity[targets.BzeroTarget], nil, bzeroTargets...)
}

// FindNBzeroTargetsOrSkipAsPolicyTarget lists the Bzero targets in the
// BastionZero organization and sets bzeroTargets to the first n Bzero targets
// found. If there are less than n Bzero targets, then the current test is
// skipped.
func FindNBzeroTargetsOrSkipAsPolicyTarget(t *testing.T, bzeroTargets ...*policies.Target) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.BzeroTarget, *http.Response, error) {
		return client.Targets.ListBzeroTargets(ctx)
	}, func(t targets.BzeroTarget) policies.Target {
		return policies.Target{ID: t.ID, Type: t.GetTargetType()}
	}, nil, bzeroTargets...)
}

// FindNClusterTargetsOrSkip lists the Cluster targets in the BastionZero
// organization and sets clusterTargets to the first n Cluster targets found. If
// there are less than n Cluster targets, then the current test is skipped.
//
// If you need the targets mapped as the policy type (policies.Cluster), use
// FindNClusterTargetsOrSkipAsPolicyCluster() instead.
func FindNClusterTargetsOrSkip(t *testing.T, clusterTargets ...*targets.ClusterTarget) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.ClusterTarget, *http.Response, error) {
		return client.Targets.ListClusterTargets(ctx)
	}, identity[targets.ClusterTarget], nil, clusterTargets...)
}

// FindNClusterTargetsOrSkipAsPolicyCluster lists the Cluster targets in the BastionZero
// organization and sets clusterTargets to the first n Cluster targets found. If
// there are less than n Cluster targets, then the current test is skipped.
func FindNClusterTargetsOrSkipAsPolicyCluster(t *testing.T, clusterTargets ...*policies.Cluster) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.ClusterTarget, *http.Response, error) {
		return client.Targets.ListClusterTargets(ctx)
	}, func(t targets.ClusterTarget) policies.Cluster {
		return policies.Cluster{ID: t.ID}
	}, nil, clusterTargets...)
}

// FindNDACTargetsOrSkip lists the DAC targets in the BastionZero organization
// and sets dacTargets to the first n DAC targets found. If there are less than
// n DAC targets, then the current test is skipped.
func FindNDACTargetsOrSkip(t *testing.T, dacTargets ...*targets.DynamicAccessConfiguration) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.DynamicAccessConfiguration, *http.Response, error) {
		return client.Targets.ListDynamicAccessConfigurations(ctx)
	}, identity[targets.DynamicAccessConfiguration], nil, dacTargets...)
}

// FindNDbTargetsOrSkip lists the Db targets in the BastionZero organization and
// sets dbTargets to the first n Db targets found. If there are less than n Db
// targets, then the current test is skipped.
//
// If you need the targets mapped as the policy type (policies.Target), use
// FindNDbTargetsOrSkipAsPolicyTarget() instead.
func FindNDbTargetsOrSkip(t *testing.T, dbTargets ...*targets.DatabaseTarget) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.DatabaseTarget, *http.Response, error) {
		return client.Targets.ListDatabaseTargets(ctx)
	}, identity[targets.DatabaseTarget], func(t targets.DatabaseTarget) bool {
		// IMPORTANT: We must filter out targets that are concurrently being
		// created by other parallel acceptance tests because they could be
		// deleted by the time the caller of this function uses them
		return !strings.HasPrefix(t.Name, TestNamePrefix)
	}, dbTargets...)
}

// FindNDbTargetsOrSkipAsPolicyTarget lists the Db targets in the BastionZero
// organization and sets dbTargets to the first n Db targets found. If there are
// less than n Db targets, then the current test is skipped.
func FindNDbTargetsOrSkipAsPolicyTarget(t *testing.T, dbTargets ...*policies.Target) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.DatabaseTarget, *http.Response, error) {
		return client.Targets.ListDatabaseTargets(ctx)
	}, func(t targets.DatabaseTarget) policies.Target {
		return policies.Target{ID: t.ID, Type: t.GetTargetType()}
	}, func(t targets.DatabaseTarget) bool {
		// IMPORTANT: We must filter out targets that are concurrently being
		// created by other parallel acceptance tests because they could be
		// deleted by the time the caller of this function uses them
		return !strings.HasPrefix(t.Name, TestNamePrefix)
	}, dbTargets...)
}

// FindNWebTargetsOrSkip lists the Web targets in the BastionZero organization
// and sets webTargets to the first n Web targets found. If there are less than
// n Web targets, then the current test is skipped.
func FindNWebTargetsOrSkip(t *testing.T, webTargets ...*targets.WebTarget) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]targets.WebTarget, *http.Response, error) {
		return client.Targets.ListWebTargets(ctx)
	}, identity[targets.WebTarget], nil, webTargets...)
}

// FindNTargetConnectPoliciesOrSkip lists the target connect policies in the
// BastionZero organization and sets targetConnectPolicies to the first n target
// connect policies found. If there are less than n target connect policies,
// then the current test is skipped.
func FindNTargetConnectPoliciesOrSkip(t *testing.T, targetConnectPolicies ...*policies.TargetConnectPolicy) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]policies.TargetConnectPolicy, *http.Response, error) {
		return client.Policies.ListTargetConnectPolicies(ctx, nil)
	}, identity[policies.TargetConnectPolicy], func(p policies.TargetConnectPolicy) bool {
		// IMPORTANT: We must filter out policies that are concurrently being
		// created by other parallel acceptance tests because they could be
		// deleted by the time the caller of this function uses them
		return !strings.HasPrefix(p.Name, TestNamePrefix)
	}, targetConnectPolicies...)
}

// FindNKubernetesPoliciesOrSkip lists the Kubernetes policies in the
// BastionZero organization and sets kubernetesPolicies to the first n
// Kubernetes policies found. If there are less than n Kubernetes policies, then
// the current test is skipped.
func FindNKubernetesPoliciesOrSkip(t *testing.T, kubernetesPolicies ...*policies.KubernetesPolicy) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]policies.KubernetesPolicy, *http.Response, error) {
		return client.Policies.ListKubernetesPolicies(ctx, nil)
	}, identity[policies.KubernetesPolicy], func(p policies.KubernetesPolicy) bool {
		// IMPORTANT: We must filter out policies that are concurrently being
		// created by other parallel acceptance tests because they could be
		// deleted by the time the caller of this function uses them
		return !strings.HasPrefix(p.Name, TestNamePrefix)
	}, kubernetesPolicies...)
}

// FindNProxyPoliciesOrSkip lists the proxy policies in the BastionZero
// organization and sets proxyPolicies to the first n proxy policies found. If
// there are less than n proxy policies, then the current test is skipped.
func FindNProxyPoliciesOrSkip(t *testing.T, proxyPolicies ...*policies.ProxyPolicy) {
	FindNAPIObjectsOrSkip(t, func(client *bzapi.Client, ctx context.Context) ([]policies.ProxyPolicy, *http.Response, error) {
		return client.Policies.ListProxyPolicies(ctx, nil)
	}, identity[policies.ProxyPolicy], func(p policies.ProxyPolicy) bool {
		// IMPORTANT: We must filter out policies that are concurrently being
		// created by other parallel acceptance tests because they could be
		// deleted by the time the caller of this function uses them
		return !strings.HasPrefix(p.Name, TestNamePrefix)
	}, proxyPolicies...)
}

func ToTerraformStringList(arr []string) string {
	// Source: https://stackoverflow.com/questions/24489384/how-to-print-the-values-of-slices#comment126502244_53672500
	return strings.ReplaceAll(fmt.Sprintf("%+q", arr), "\" \"", "\",\"")
}

func ExpandValuesCheckMapToSingleCheck[T any](resourceName string, apiObject *T, getValuesCheckMapFunc func(apiObject *T) map[string]string) resource.TestCheckFunc {
	valuesCheckMap := getValuesCheckMapFunc(apiObject)
	var checkFuncs []resource.TestCheckFunc
	for attr, value := range valuesCheckMap {
		if value != "" {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(resourceName, attr, value))
		} else {
			// "" denotes attribute should be unset (null)
			checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(resourceName, attr))
		}
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func SafePrettyInt(value *int) string {
	if value == nil {
		return "<nil>"
	} else {
		return strconv.Itoa(*value)
	}
}

// TerraformObjectToString returns a human-readable representation of the Object
// value that can be used in Terraform configurations. This function is
// equivalent to basetypes.Object.String() except that null object attributes
// are omitted from the final string.
func TerraformObjectToString(o types.Object) string {
	// Source:
	// https://github.com/hashicorp/terraform-plugin-framework/blob/3b275fbb481bd598b3e020f2560cddc70b885098/types/basetypes/object_value.go#L364-L395
	//
	// This function is equivalent to basetypes.Object.String() except that null
	// object fields are omitted

	if o.IsUnknown() {
		return attr.UnknownValueString
	}

	if o.IsNull() {
		return attr.NullValueString
	}

	// We want the output to be consistent, so we sort the output by key
	keys := make([]string, 0, len(o.Attributes()))
	for k := range o.Attributes() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res strings.Builder

	res.WriteString("{")
	for i, k := range keys {
		// Omit null attributes from final string representation
		if o.Attributes()[k].IsNull() {
			continue
		}
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(fmt.Sprintf(`"%s":%s`, k, o.Attributes()[k].String()))
	}
	res.WriteString("}")

	return res.String()
}
