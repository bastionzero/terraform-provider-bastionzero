package kubernetes_test

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
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccKubernetesPolicy_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTargetConnectPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccTargetConnectPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetConnectPolicyExists(resourceName, &policy),
					testAccCheckTargetConnectPolicyAttributes(t, &policy, &expectedTargetConnectPolicy{
						Name:          &rName,
						Description:   bastionzero.PtrTo(""),
						Subjects:      &[]policies.Subject{},
						Groups:        &[]policies.Group{},
						Environments:  &[]policies.Environment{},
						Clusters:      &[]policies.Cluster{},
						ClusterUsers:  &[]policies.ClusterUser{},
						ClusterGroups: &[]policies.ClusterGroup{},
					}),
					testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					// Check that unspecified values remain null
					resource.TestCheckNoResourceAttr(resourceName, "subjects"),
					resource.TestCheckNoResourceAttr(resourceName, "groups"),
					resource.TestCheckNoResourceAttr(resourceName, "environments"),
					resource.TestCheckNoResourceAttr(resourceName, "clusters"),
					resource.TestCheckNoResourceAttr(resourceName, "cluster_users"),
					resource.TestCheckNoResourceAttr(resourceName, "cluster_groups"),
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

func testAccTargetConnectPolicyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  name = %[1]q
}
`, rName)
}

func testAccTargetConnectPolicyConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  description = %[2]q
  name = %[1]q
}
`, rName, description)
}

func testAccTargetConnectPolicyConfigSubjects(rName string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  subjects = %[2]s
  name = %[1]q
}
`, rName, subjects.String())
}

func testAccTargetConnectPolicyConfigGroups(rName string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  groups = %[2]s
  name = %[1]q
}
`, rName, groups.String())
}

func testAccTargetConnectPolicyConfigEnvironments(rName string, environments []string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  environments = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(environments))
}

func testAccKubernetesPolicyConfigClusters(rName string, clusters []string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  clusters = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(clusters))
}

func testAccKubernetesPolicyConfigClusterUsers(rName string, clusterUsers []string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  cluster_users = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(clusterUsers))
}

func testAccKubernetesPolicyConfigClusterGroups(rName string, clusterGroups []string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  cluster_groups = %[2]s
  name = %[1]q
}
`, rName, acctest.ToTerraformStringList(clusterGroups))
}

type expectedTargetConnectPolicy struct {
	Name          *string
	Description   *string
	Subjects      *[]policies.Subject
	Groups        *[]policies.Group
	Environments  *[]policies.Environment
	Clusters      *[]policies.Cluster
	ClusterUsers  *[]policies.ClusterUser
	ClusterGroups *[]policies.ClusterGroup
}

func testAccCheckTargetConnectPolicyAttributes(t *testing.T, policy *policies.KubernetesPolicy, expected *expectedTargetConnectPolicy) resource.TestCheckFunc {
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
		if expected.Clusters != nil && !assert.ElementsMatch(t, *expected.Clusters, policy.GetClusters()) {
			return fmt.Errorf("Bad clusters, expected \"%s\", got: %#v", *expected.Clusters, policy.Clusters)
		}
		if expected.ClusterUsers != nil && !assert.ElementsMatch(t, *expected.ClusterUsers, policy.GetClusterUsers()) {
			return fmt.Errorf("Bad cluster_users, expected \"%s\", got: %#v", *expected.ClusterUsers, policy.ClusterUsers)
		}
		if expected.ClusterGroups != nil && !assert.ElementsMatch(t, *expected.ClusterGroups, policy.GetClusterGroups()) {
			return fmt.Errorf("Bad cluster_groups, expected \"%s\", got: %#v", *expected.ClusterGroups, policy.ClusterGroups)
		}

		return nil
	}
}

func testAccCheckTargetConnectPolicyExists(namedTFResource string, policy *policies.KubernetesPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.KubernetesPolicy, *http.Response, error) {
		return c.Policies.GetKubernetesPolicy(ctx, id)
	})
}

func testAccCheckResourceTargetConnectPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.Kubernetes)),
	)
}

func testAccCheckTargetConnectPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bastionzero_kubernetes_policy" {
			continue
		}

		// Try to find the policy
		_, _, err := acctest.APIClient.Policies.GetKubernetesPolicy(context.Background(), rs.Primary.ID)
		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
			return fmt.Errorf("Error waiting for Kubernetes policy (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
