package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/proxy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen"
	"pgregory.net/rapid"
)

func ProxyPolicySchemaGen(ctx context.Context) *rapid.Generator[proxy.ProxyPolicyModel] {
	return rapid.Custom(func(t *rapid.T) proxy.ProxyPolicyModel {
		return proxy.ProxyPolicyModel{
			ID:           tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:         tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:         tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Proxy)).Draw(t, "TypeOrNull"),
			Description:  tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:     tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:       tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments: tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Targets:      tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetsGen()).Draw(t, "TargetsOrNull"),
			TargetUsers:  tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetUsersGen()).Draw(t, "TargetUsersOrNull"),
		}
	})
}
