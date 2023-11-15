package policy_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestFlatExpandPolicySubjects_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicySubjectGen()).Draw(t, "Subjects")

		// Flatten the generated BastionZero type into a TF type
		flattened := policy.FlattenPolicySubjects(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := policy.ExpandPolicySubjects(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandPolicyGroups_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyGroupGen()).Draw(t, "Groups")

		// Flatten the generated BastionZero type into a TF type
		flattened := policy.FlattenPolicyGroups(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := policy.ExpandPolicyGroups(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandPolicyEnvironments_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyEnvironmentGen()).Draw(t, "Environments")

		// Flatten the generated BastionZero type into a TF type
		flattened := policy.FlattenPolicyEnvironments(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := policy.ExpandPolicyEnvironments(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandPolicyTargets_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyTargetGen()).Draw(t, "Targets")

		// Flatten the generated BastionZero type into a TF type
		flattened := policy.FlattenPolicyTargets(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := policy.ExpandPolicyTargets(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}

func TestFlatExpandPolicyTargetUsers_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.PolicyTargetUserGen()).Draw(t, "TargetUsers")

		// Flatten the generated BastionZero type into a TF type
		flattened := policy.FlattenPolicyTargetUsers(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := policy.ExpandPolicyTargetUsers(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, genAPI, expanded)
	})
}
