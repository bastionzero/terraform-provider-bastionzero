package targetconnect_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/targetconnect"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfpolicygen"
	"github.com/stretchr/testify/require"

	"pgregory.net/rapid"
)

func TestFlatExpandPolicyVerbs_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyVerbGen()).Draw(t, "Verbs")

		// Flatten the generated BastionZero type into a TF type
		flattened := targetconnect.FlattenPolicyVerbs(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := targetconnect.ExpandPolicyVerbs(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandTargetConnect_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genBzTargetConnectPolicy := bzpolicygen.TargetConnectPolicyGen().Draw(t, "api")
		genTFSchemaTargetConnectPolicy := tfpolicygen.TargetConnectPolicySchemaGen(context.Background()).Draw(t, "schema")

		// Flatten generated BastionZero API policy into Terraform schema type
		targetconnect.SetTargetConnectPolicyAttributes(context.Background(), &genTFSchemaTargetConnectPolicy, &genBzTargetConnectPolicy, false)

		// Then expand the flattened value back to a BastionZero API policy type
		expanded := targetconnect.ExpandTargetConnectPolicy(context.Background(), &genTFSchemaTargetConnectPolicy)

		// And assert no data loss occurred when converting by asserting the
		// expanded value matches the original, generated policy
		require.EqualValues(t, genBzTargetConnectPolicy.Name, expanded.Name)
		require.EqualValues(t, genBzTargetConnectPolicy.GetPolicyType(), expanded.GetPolicyType())
		require.EqualValues(t, genBzTargetConnectPolicy.Description, expanded.Description)
		require.EqualValues(t, genBzTargetConnectPolicy.Subjects, expanded.Subjects)
		require.EqualValues(t, genBzTargetConnectPolicy.Groups, expanded.Groups)
		require.EqualValues(t, genBzTargetConnectPolicy.Environments, expanded.Environments)
		require.EqualValues(t, genBzTargetConnectPolicy.Targets, expanded.Targets)
		require.EqualValues(t, genBzTargetConnectPolicy.TargetUsers, expanded.TargetUsers)
		require.EqualValues(t, genBzTargetConnectPolicy.Verbs, expanded.Verbs)
	})
}
