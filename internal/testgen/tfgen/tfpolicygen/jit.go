package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/jit"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen"
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
			ID:                    tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:                  tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:                  tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description:           tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:              tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:                tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			ChildPolicies:         tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaChildPoliciesGen()).Draw(t, "ChildPoliciesOrNull"),
			AutomaticallyApproved: tfgen.BoolWithValueOrNullGen(ctx).Draw(t, "AutomaticallyApprovedOrNull"),
			Duration:              tfgen.Int64WithValueOrNullGen(ctx).Draw(t, "DurationOrNull"),
		}
	})
}
