package kubernetes

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// kubernetesPolicyModel maps the kubernetes policy schema data.
type kubernetesPolicyModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	Subjects      types.Set    `tfsdk:"subjects"`
	Groups        types.Set    `tfsdk:"groups"`
	Environments  types.Set    `tfsdk:"environments"`
	Clusters      types.Set    `tfsdk:"clusters"`
	ClusterUsers  types.Set    `tfsdk:"cluster_users"`
	ClusterGroups types.Set    `tfsdk:"cluster_groups"`
}

func (m *kubernetesPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *kubernetesPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *kubernetesPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *kubernetesPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *kubernetesPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *kubernetesPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *kubernetesPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *kubernetesPolicyModel) GetGroups() types.Set   { return m.Groups }

// setKubernetesPolicyAttributes populates the TF schema data from a kubernetes
// policy
func setKubernetesPolicyAttributes(ctx context.Context, schema *kubernetesPolicyModel, apiPolicy *policies.KubernetesPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)

	// See comment in SetBasePolicyAttributes that explains this conditional
	// logic
	if !schema.Environments.IsNull() || len(apiPolicy.GetEnvironments()) != 0 || modelIsDataSource {
		schema.Environments = policy.FlattenPolicyEnvironments(ctx, apiPolicy.GetEnvironments())
	}
	if !schema.Clusters.IsNull() || len(apiPolicy.GetClusters()) != 0 || modelIsDataSource {
		schema.Clusters = FlattenPolicyClusters(ctx, apiPolicy.GetClusters())
	}
	if !schema.ClusterUsers.IsNull() || len(apiPolicy.GetClusterUsers()) != 0 || modelIsDataSource {
		schema.ClusterUsers = FlattenPolicyClusterUsers(ctx, apiPolicy.GetClusterUsers())
	}
	if !schema.ClusterGroups.IsNull() || len(apiPolicy.GetClusterGroups()) != 0 || modelIsDataSource {
		schema.ClusterGroups = FlattenPolicyClusterGroups(ctx, apiPolicy.GetClusterGroups())
	}
}

func ExpandPolicyClusters(ctx context.Context, tfSet types.Set) []policies.Cluster {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.Cluster {
		return policies.Cluster{ID: m}
	})
}

func FlattenPolicyClusters(ctx context.Context, apiObject []policies.Cluster) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.Cluster) attr.Value {
		return types.StringValue(m.ID)
	})
}

func ExpandPolicyClusterUsers(ctx context.Context, tfSet types.Set) []policies.ClusterUser {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.ClusterUser {
		return policies.ClusterUser{Name: m}
	})
}

func FlattenPolicyClusterUsers(ctx context.Context, apiObject []policies.ClusterUser) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.ClusterUser) attr.Value {
		return types.StringValue(m.Name)
	})
}

func ExpandPolicyClusterGroups(ctx context.Context, tfSet types.Set) []policies.ClusterGroup {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.ClusterGroup {
		return policies.ClusterGroup{Name: m}
	})
}

func FlattenPolicyClusterGroups(ctx context.Context, apiObject []policies.ClusterGroup) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.ClusterGroup) attr.Value {
		return types.StringValue(m.Name)
	})
}

func makeKubernetesPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.Kubernetes)
	attributes["environments"] = policy.PolicyEnvironmentsAttribute()
	attributes["clusters"] = schema.SetAttribute{
		Description: "Set of Cluster target ID(s) that this policy applies to.",
		ElementType: types.StringType,
		Optional:    true,
	}
	attributes["cluster_users"] = schema.SetAttribute{
		Description: "Set of Kubernetes RBAC subject usernames that this policy applies to. " +
			"See https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-subjects for more details.",
		ElementType: types.StringType,
		Optional:    true,
	}
	attributes["cluster_groups"] = schema.SetAttribute{
		Description: "Set of Kubernetes RBAC subject group names that this policy applies to. " +
			"See https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-subjects for more details.",
		ElementType: types.StringType,
		Optional:    true,
	}

	return attributes
}
