package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/targetconnect"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"pgregory.net/rapid"
)

func PolicySchemaVerbsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return targetconnect.FlattenPolicyVerbs(context.Background(), rapid.SliceOf(bzpolicygen.PolicyVerbGen()).Draw(t, "SchemaVerbs"))
	})
}

func TargetConnectPolicySchemaGen(ctx context.Context) *rapid.Generator[targetconnect.TargetConnectPolicyModel] {
	return rapid.Custom(func(t *rapid.T) targetconnect.TargetConnectPolicyModel {
		return targetconnect.TargetConnectPolicyModel{
			ID:           tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:         tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:         tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.TargetConnect)).Draw(t, "TypeOrNull"),
			Description:  tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:     tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:       tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments: tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Targets:      tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetsGen()).Draw(t, "TargetsOrNull"),
			TargetUsers:  tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetUsersGen()).Draw(t, "TargetUsersOrNull"),
			Verbs:        tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaVerbsGen()).Draw(t, "VerbsOrNull"),
		}
	})
}
