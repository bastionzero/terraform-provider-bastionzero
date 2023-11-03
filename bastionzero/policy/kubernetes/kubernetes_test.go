package kubernetes_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/kubernetes"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen/tfpolicygen"
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
		require.EqualValues(t, genKubernetesPolicy.GetDescription(), expanded.GetDescription())
		require.EqualValues(t, genKubernetesPolicy.GetSubjects(), expanded.GetSubjects())
		require.EqualValues(t, genKubernetesPolicy.GetGroups(), expanded.GetGroups())
		require.EqualValues(t, genKubernetesPolicy.GetEnvironments(), expanded.GetEnvironments())
		require.EqualValues(t, genKubernetesPolicy.GetClusters(), expanded.GetClusters())
		require.EqualValues(t, genKubernetesPolicy.GetClusterUsers(), expanded.GetClusterUsers())
		require.EqualValues(t, genKubernetesPolicy.GetClusterGroups(), expanded.GetClusterGroups())
	})
}
