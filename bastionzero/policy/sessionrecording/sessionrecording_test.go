package sessionrecording_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/sessionrecording"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfpolicygen"
	"github.com/stretchr/testify/require"

	"pgregory.net/rapid"
)

func TestFlatExpandSessionRecording_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genSessionRecordingPolicy := bzpolicygen.SessionRecordingPolicyGen().Draw(t, "api")
		genTFSchemaSessionRecordingPolicy := tfpolicygen.SessionRecordingPolicySchemaGen(context.Background()).Draw(t, "schema")

		// Flatten generated BastionZero API policy into Terraform schema type
		sessionrecording.SetSessionRecordingPolicyAttributes(context.Background(), &genTFSchemaSessionRecordingPolicy, &genSessionRecordingPolicy, false)

		// Then expand the flattened value back to a BastionZero API policy type
		expanded := sessionrecording.ExpandSessionRecordingPolicy(context.Background(), &genTFSchemaSessionRecordingPolicy)

		// And assert no data loss occurred when converting by asserting the
		// expanded value matches the original, generated policy
		require.EqualValues(t, genSessionRecordingPolicy.Name, expanded.Name)
		require.EqualValues(t, genSessionRecordingPolicy.GetPolicyType(), expanded.GetPolicyType())
		require.EqualValues(t, genSessionRecordingPolicy.GetDescription(), expanded.GetDescription())
		require.EqualValues(t, genSessionRecordingPolicy.GetSubjects(), expanded.GetSubjects())
		require.EqualValues(t, genSessionRecordingPolicy.GetGroups(), expanded.GetGroups())
		require.EqualValues(t, genSessionRecordingPolicy.GetRecordInput(), expanded.GetRecordInput())
	})
}
