package kubernetes_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/kubernetes"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfpolicygen"
	"github.com/stretchr/testify/require"

	"pgregory.net/rapid"
)

func TestFlatExpandClusters_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyClusterGen()).Draw(t, "Clusters")

		// Flatten the generated BastionZero type into a TF type
		flattened := kubernetes.FlattenPolicyClusters(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := kubernetes.ExpandPolicyClusters(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandClusterUsers_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyClusterUserGen()).Draw(t, "ClusterUsers")

		// Flatten the generated BastionZero type into a TF type
		flattened := kubernetes.FlattenPolicyClusterUsers(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := kubernetes.ExpandPolicyClusterUsers(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandClusterGroups_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyClusterGroupGen()).Draw(t, "ClusterGroups")

		// Flatten the generated BastionZero type into a TF type
		flattened := kubernetes.FlattenPolicyClusterGroups(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := kubernetes.ExpandPolicyClusterGroups(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandKubernetes_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genKubernetesPolicy := bzpolicygen.KubernetesPolicyGen().Draw(t, "api")
		genTFSchemaKubernetesPolicy := tfpolicygen.KubernetesPolicySchemaGen(context.Background()).Draw(t, "schema")

		// Flatten generated BastionZero API policy into Terraform schema type
		kubernetes.SetKubernetesPolicyAttributes(context.Background(), &genTFSchemaKubernetesPolicy, &genKubernetesPolicy, false)

		// Then expand the flattened value back to a BastionZero API policy type
		expanded := kubernetes.ExpandKubernetesPolicy(context.Background(), &genTFSchemaKubernetesPolicy)

		// And assert no data loss occurred when converting by asserting the
		// expanded value matches the original, generated policy
		require.EqualValues(t, genKubernetesPolicy.Name, expanded.Name)
		require.EqualValues(t, genKubernetesPolicy.GetPolicyType(), expanded.GetPolicyType())
		require.EqualValues(t, genKubernetesPolicy.Description, expanded.Description)
		require.EqualValues(t, genKubernetesPolicy.Subjects, expanded.Subjects)
		require.EqualValues(t, genKubernetesPolicy.Groups, expanded.Groups)
		require.EqualValues(t, genKubernetesPolicy.Environments, expanded.Environments)
		require.EqualValues(t, genKubernetesPolicy.Clusters, expanded.Clusters)
		require.EqualValues(t, genKubernetesPolicy.ClusterUsers, expanded.ClusterUsers)
		require.EqualValues(t, genKubernetesPolicy.ClusterGroups, expanded.ClusterGroups)
	})
}
