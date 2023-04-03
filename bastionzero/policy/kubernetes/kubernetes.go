package kubernetes

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KubernetesPolicyModel maps the kubernetes policy schema data.
type KubernetesPolicyModel struct {
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

func (m *KubernetesPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *KubernetesPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *KubernetesPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *KubernetesPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *KubernetesPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *KubernetesPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *KubernetesPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *KubernetesPolicyModel) GetGroups() types.Set   { return m.Groups }

// SetKubernetesPolicyAttributes populates the TF schema data from a kubernetes
// policy
func SetKubernetesPolicyAttributes(ctx context.Context, schema *KubernetesPolicyModel, apiPolicy *policies.KubernetesPolicy, modelIsDataSource bool) {
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
	// if !schema.ClusterGroups.IsNull() || len(apiPolicy.GetClusterGroups()) != 0 || modelIsDataSource {
	// 	schema.ClusterGroups = FlattenPolicyClusterGroups(ctx, apiPolicy.GetClusterGroups())
	// }
	schema.ClusterGroups = FlattenPolicyClusterGroups(ctx, apiPolicy.GetClusterGroups())
}

func ExpandKubernetesPolicy(ctx context.Context, schema *KubernetesPolicyModel) *policies.KubernetesPolicy {
	p := new(policies.KubernetesPolicy)
	p.Name = schema.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, schema.Description)
	p.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, schema.Subjects))
	p.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, schema.Groups))
	p.Environments = bastionzero.PtrTo(policy.ExpandPolicyEnvironments(ctx, schema.Environments))
	p.Clusters = bastionzero.PtrTo(ExpandPolicyClusters(ctx, schema.Clusters))
	p.ClusterUsers = bastionzero.PtrTo(ExpandPolicyClusterUsers(ctx, schema.ClusterUsers))
	p.ClusterGroups = bastionzero.PtrTo(ExpandPolicyClusterGroups(ctx, schema.ClusterGroups))

	return p
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
