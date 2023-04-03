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
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
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
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccKubernetesPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:          &rName,
						Description:   bastionzero.PtrTo(""),
						Subjects:      &[]policies.Subject{},
						Groups:        &[]policies.Group{},
						Environments:  &[]policies.Environment{},
						Clusters:      &[]policies.Cluster{},
						ClusterUsers:  &[]policies.ClusterUser{},
						ClusterGroups: &[]policies.ClusterGroup{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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

func TestAccKubernetesPolicy_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Policies.DeleteKubernetesPolicy(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKubernetesPolicy_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigBasic(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name: &rName1,
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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
				Config: testAccKubernetesPolicyConfigBasic(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name: &rName2,
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_Description(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy
	desc1 := "desc1"
	desc2 := "desc2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigDescription(rName, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:        &rName,
						Description: &desc1,
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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
				Config: testAccKubernetesPolicyConfigDescription(rName, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:        &rName,
						Description: &desc2,
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Verify setting to empty string clears
			{
				Config: testAccKubernetesPolicyConfigDescription(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:        &rName,
						Description: bastionzero.PtrTo(""),
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_Subjects(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var p policies.KubernetesPolicy
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
	acctest.FindNUsersOrSkip(t, subject1, subject2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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
				Config: testAccKubernetesPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{*subject1, *subject2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{*subject1, *subject2},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject1.ID, "type": string(subject1.Type)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "subjects.*", map[string]string{"id": subject2.ID, "type": string(subject2.Type)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigSubjects(rName, policy.FlattenPolicySubjects(ctx, []policies.Subject{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Subjects: &[]policies.Subject{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "subjects.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_Groups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var p policies.KubernetesPolicy
	group1 := new(policies.Group)
	group2 := new(policies.Group)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNGroupsOrSkip(t, group1, group2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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
				Config: testAccKubernetesPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{*group1, *group2})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{*group1, *group2},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group1.ID, "name": string(group1.Name)}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "groups.*", map[string]string{"id": group2.ID, "name": string(group2.Name)}),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigGroups(rName, policy.FlattenPolicyGroups(ctx, []policies.Group{})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:   &rName,
						Groups: &[]policies.Group{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_Environments(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var p policies.KubernetesPolicy
	env1 := new(policies.Environment)
	env2 := new(policies.Environment)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNEnvironmentsOrSkip(t, env1, env2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigEnvironments(rName, []string{env1.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
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
				Config: testAccKubernetesPolicyConfigEnvironments(rName, []string{env1.ID, env2.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{*env1, *env2},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env1.ID),
					resource.TestCheckTypeSetElemAttr(resourceName, "environments.*", env2.ID),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigEnvironments(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:         &rName,
						Environments: &[]policies.Environment{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "environments.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_Clusters(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var p policies.KubernetesPolicy
	target1 := new(policies.Cluster)
	target2 := new(policies.Cluster)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNClusterTargetsOrSkip(t, target1, target2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigClusters(rName, []string{target1.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Clusters: &[]policies.Cluster{*target1},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "clusters.*", target1.ID),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update clusters
			{
				Config: testAccKubernetesPolicyConfigClusters(rName, []string{target1.ID, target2.ID}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Clusters: &[]policies.Cluster{*target1, *target2},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "clusters.*", target1.ID),
					resource.TestCheckTypeSetElemAttr(resourceName, "clusters.*", target2.ID),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigClusters(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &p),
					testAccCheckKubernetesPolicyAttributes(t, &p, &expectedKubernetesPolicy{
						Name:     &rName,
						Clusters: &[]policies.Cluster{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_ClusterUsers(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigClusterUsers(rName, []string{"foo"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:         &rName,
						ClusterUsers: &[]policies.ClusterUser{{Name: "foo"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_users.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_users.*", "foo"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update cluster users
			{
				Config: testAccKubernetesPolicyConfigClusterUsers(rName, []string{"bar"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:         &rName,
						ClusterUsers: &[]policies.ClusterUser{{Name: "bar"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_users.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_users.*", "bar"),
				),
			},
			// Add another cluster user
			{
				Config: testAccKubernetesPolicyConfigClusterUsers(rName, []string{"bar", "baz"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:         &rName,
						ClusterUsers: &[]policies.ClusterUser{{Name: "bar"}, {Name: "baz"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_users.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_users.*", "bar"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_users.*", "baz"),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigClusterUsers(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:         &rName,
						ClusterUsers: &[]policies.ClusterUser{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "cluster_users.#", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesPolicy_ClusterGroups(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_kubernetes_policy.test"
	var policy policies.KubernetesPolicy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKubernetesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesPolicyConfigClusterGroups(rName, []string{"foo"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:          &rName,
						ClusterGroups: &[]policies.ClusterGroup{{Name: "foo"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_groups.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_groups.*", "foo"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update cluster groups
			{
				Config: testAccKubernetesPolicyConfigClusterGroups(rName, []string{"bar"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:          &rName,
						ClusterGroups: &[]policies.ClusterGroup{{Name: "bar"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_groups.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_groups.*", "bar"),
				),
			},
			// Add another cluster group
			{
				Config: testAccKubernetesPolicyConfigClusterGroups(rName, []string{"bar", "baz"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:          &rName,
						ClusterGroups: &[]policies.ClusterGroup{{Name: "bar"}, {Name: "baz"}},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_groups.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_groups.*", "bar"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cluster_groups.*", "baz"),
				),
			},
			// Verify setting to empty list clears
			{
				Config: testAccKubernetesPolicyConfigClusterGroups(rName, []string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesPolicyExists(resourceName, &policy),
					testAccCheckKubernetesPolicyAttributes(t, &policy, &expectedKubernetesPolicy{
						Name:          &rName,
						ClusterGroups: &[]policies.ClusterGroup{},
					}),
					testAccCheckResourceKubernetesPolicyComputedAttr(resourceName),
					// Explicit empty list in config should result in a config
					// with 0 elements (not null)
					resource.TestCheckResourceAttr(resourceName, "cluster_groups.#", "0"),
				),
			},
		},
	})
}

func testAccKubernetesPolicyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  name = %[1]q
}
`, rName)
}

func testAccKubernetesPolicyConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  description = %[2]q
  name = %[1]q
}
`, rName, description)
}

func testAccKubernetesPolicyConfigSubjects(rName string, subjects types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  subjects = %[2]s
  name = %[1]q
}
`, rName, subjects.String())
}

func testAccKubernetesPolicyConfigGroups(rName string, groups types.Set) string {
	return fmt.Sprintf(`
resource "bastionzero_kubernetes_policy" "test" {
  groups = %[2]s
  name = %[1]q
}
`, rName, groups.String())
}

func testAccKubernetesPolicyConfigEnvironments(rName string, environments []string) string {
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

type expectedKubernetesPolicy struct {
	Name          *string
	Description   *string
	Subjects      *[]policies.Subject
	Groups        *[]policies.Group
	Environments  *[]policies.Environment
	Clusters      *[]policies.Cluster
	ClusterUsers  *[]policies.ClusterUser
	ClusterGroups *[]policies.ClusterGroup
}

func testAccCheckKubernetesPolicyAttributes(t *testing.T, policy *policies.KubernetesPolicy, expected *expectedKubernetesPolicy) resource.TestCheckFunc {
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

func testAccCheckKubernetesPolicyExists(namedTFResource string, policy *policies.KubernetesPolicy) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, policy, func(c *bastionzero.Client, ctx context.Context, id string) (*policies.KubernetesPolicy, *http.Response, error) {
		return c.Policies.GetKubernetesPolicy(ctx, id)
	})
}

func testAccCheckResourceKubernetesPolicyComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckResourceAttr(resourceName, "type", string(policytype.Kubernetes)),
	)
}

func testAccCheckKubernetesPolicyDestroy(s *terraform.State) error {
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
