package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/targetconnect"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tftypesgen"
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
			ID:           tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:         tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:         tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.TargetConnect)).Draw(t, "TypeOrNull"),
			Description:  tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:     tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:       tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments: tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Targets:      tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetsGen()).Draw(t, "TargetsOrNull"),
			TargetUsers:  tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetUsersGen()).Draw(t, "TargetUsersOrNull"),
			Verbs:        tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaVerbsGen()).Draw(t, "VerbsOrNull"),
		}
	})
}
