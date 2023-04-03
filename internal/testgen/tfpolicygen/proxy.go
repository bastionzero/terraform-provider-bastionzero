package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/proxy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tftypesgen"
	"pgregory.net/rapid"
)

func ProxyPolicySchemaGen(ctx context.Context) *rapid.Generator[proxy.ProxyPolicyModel] {
	return rapid.Custom(func(t *rapid.T) proxy.ProxyPolicyModel {
		return proxy.ProxyPolicyModel{
			ID:           tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:         tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:         tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Proxy)).Draw(t, "TypeOrNull"),
			Description:  tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:     tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:       tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments: tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Targets:      tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetsGen()).Draw(t, "TargetsOrNull"),
			TargetUsers:  tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaTargetUsersGen()).Draw(t, "TargetUsersOrNull"),
		}
	})
}
