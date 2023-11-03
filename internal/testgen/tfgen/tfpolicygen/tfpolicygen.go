// Package tfpolicygen provides property based test (PBT) generators for
// BastionZero Terraform Provider policy schema model types
package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"pgregory.net/rapid"
)

func PolicySchemaIDGen() *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return basetypes.NewStringValue(rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "SchemaID"))
	})
}

func PolicySchemaNameGen() *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return basetypes.NewStringValue(rapid.String().Draw(t, "SchemaName"))
	})
}

func PolicySchemaTypeGen(policyType policytype.PolicyType) *rapid.Generator[basetypes.StringValue] {
	return rapid.Just(basetypes.NewStringValue(string(policyType)))
}

func PolicySchemaDescriptionGen() *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return basetypes.NewStringValue(rapid.String().Draw(t, "SchemaDescription"))
	})
}

func PolicySchemaSubjectsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return policy.FlattenPolicySubjects(context.Background(), rapid.SliceOf(bzpolicygen.PolicySubjectGen()).Draw(t, "SchemaSubjects"))
	})
}

func PolicySchemaGroupsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return policy.FlattenPolicyGroups(context.Background(), rapid.SliceOf(bzpolicygen.PolicyGroupGen()).Draw(t, "SchemaGroups"))
	})
}

func PolicySchemaEnvironmentsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return policy.FlattenPolicyEnvironments(context.Background(), rapid.SliceOf(bzpolicygen.PolicyEnvironmentGen()).Draw(t, "SchemaEnvironments"))
	})
}

func PolicySchemaTargetsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return policy.FlattenPolicyTargets(context.Background(), rapid.SliceOf(bzpolicygen.PolicyTargetGen()).Draw(t, "SchemaTargets"))
	})
}

func PolicySchemaTargetUsersGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return policy.FlattenPolicyTargetUsers(context.Background(), rapid.SliceOf(bzpolicygen.PolicyTargetUserGen()).Draw(t, "SchemaTargetUsers"))
	})
}
