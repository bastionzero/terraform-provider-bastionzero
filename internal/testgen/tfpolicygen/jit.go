package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/jit"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tftypesgen"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"pgregory.net/rapid"
)

func PolicySchemaChildPoliciesGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return jit.FlattenChildPolicies(context.Background(), rapid.SliceOf(bzpolicygen.ChildPolicyGen()).Draw(t, "SchemaChildPolicies"))
	})
}

func JITPolicySchemaGen(ctx context.Context) *rapid.Generator[jit.JITPolicyModel] {
	return rapid.Custom(func(t *rapid.T) jit.JITPolicyModel {
		return jit.JITPolicyModel{
			ID:                    tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:                  tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:                  tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description:           tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:              tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:                tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			ChildPolicies:         tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaChildPoliciesGen()).Draw(t, "ChildPoliciesOrNull"),
			AutomaticallyApproved: tftypesgen.BoolWithValueOrNullGen(ctx).Draw(t, "AutomaticallyApprovedOrNull"),
			Duration:              tftypesgen.Int64WithValueOrNullGen(ctx).Draw(t, "DurationOrNull"),
		}
	})
}
