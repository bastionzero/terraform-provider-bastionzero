package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/kubernetes"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tftypesgen"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"pgregory.net/rapid"
)

func PolicySchemaClustersGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return kubernetes.FlattenPolicyClusters(context.Background(), rapid.SliceOf(bzpolicygen.PolicyClusterGen()).Draw(t, "SchemaClusters"))
	})
}

func PolicySchemaClusterUsersGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return kubernetes.FlattenPolicyClusterUsers(context.Background(), rapid.SliceOf(bzpolicygen.PolicyClusterUserGen()).Draw(t, "SchemaClusterUsers"))
	})
}

func PolicySchemaClusterGroupsGen() *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return kubernetes.FlattenPolicyClusterGroups(context.Background(), rapid.SliceOf(bzpolicygen.PolicyClusterGroupGen()).Draw(t, "SchemaClusterGroups"))
	})
}

func KubernetesPolicySchemaGen(ctx context.Context) *rapid.Generator[kubernetes.KubernetesPolicyModel] {
	return rapid.Custom(func(t *rapid.T) kubernetes.KubernetesPolicyModel {
		return kubernetes.KubernetesPolicyModel{
			ID:            tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:          tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:          tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description:   tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:      tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:        tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments:  tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Clusters:      tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClustersGen()).Draw(t, "ClustersOrNull"),
			ClusterUsers:  tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClusterUsersGen()).Draw(t, "ClusterUsersOrNull"),
			ClusterGroups: tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClusterGroupsGen()).Draw(t, "ClusterGroupsOrNull"),
		}
	})
}
