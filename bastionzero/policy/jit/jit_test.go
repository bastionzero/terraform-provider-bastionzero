package jit_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/jit"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen/tfpolicygen"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestFlatExpandChildPolicies_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := rapid.SliceOf(bzpolicygen.ChildPolicyGen()).Draw(t, "ChildPolicies")

		// Flatten the generated BastionZero type into a TF type
		flattened := jit.FlattenChildPolicies(context.Background(), genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := jit.ExpandChildPolicies(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, childPoliciesToIDs(genAPI), expanded)
	})
}

func TestFlatExpandJIT_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genJITPolicy := bzpolicygen.JITPolicyGen().Draw(t, "api")
		genTFSchemaJITPolicy := tfpolicygen.JITPolicySchemaGen(context.Background()).Draw(t, "schema")

		// Flatten generated BastionZero API policy into Terraform schema type
		jit.SetJITPolicyAttributes(context.Background(), &genTFSchemaJITPolicy, &genJITPolicy, false)

		// Then expand the flattened value back to a BastionZero API policy type
		expanded := jit.ExpandJITPolicy(context.Background(), &genTFSchemaJITPolicy)

		// And assert no data loss occurred when converting by asserting the
		// expanded value matches the original, generated policy
		require.EqualValues(t, genJITPolicy.Name, expanded.Name)
		require.EqualValues(t, genJITPolicy.GetDescription(), expanded.Description)
		require.EqualValues(t, genJITPolicy.GetSubjects(), expanded.Subjects)
		require.EqualValues(t, genJITPolicy.GetGroups(), expanded.Groups)
		require.EqualValues(t, childPoliciesToIDs(genJITPolicy.GetChildPolicies()), expanded.ChildPolicies)
		require.EqualValues(t, genJITPolicy.GetAutomaticallyApproved(), expanded.AutomaticallyApproved)
		require.EqualValues(t, genJITPolicy.GetDuration(), expanded.Duration)
	})
}

func childPoliciesToIDs(childPolicies []policies.ChildPolicy) []string {
	ids := make([]string, 0)
	for _, p := range childPolicies {
		ids = append(ids, p.ID)
	}
	return ids
}
