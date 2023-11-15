package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/kubernetes"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen"
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
			ID:            tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:          tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:          tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description:   tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:      tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:        tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			Environments:  tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaEnvironmentsGen()).Draw(t, "EnvrionmentsOrNull"),
			Clusters:      tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClustersGen()).Draw(t, "ClustersOrNull"),
			ClusterUsers:  tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClusterUsersGen()).Draw(t, "ClusterUsersOrNull"),
			ClusterGroups: tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaClusterGroupsGen()).Draw(t, "ClusterGroupsOrNull"),
		}
	})
}
