package proxy

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProxyPolicyModel maps the proxy policy schema data.
type ProxyPolicyModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	Subjects     types.Set    `tfsdk:"subjects"`
	Groups       types.Set    `tfsdk:"groups"`
	Environments types.Set    `tfsdk:"environments"`
	Targets      types.Set    `tfsdk:"targets"`
	TargetUsers  types.Set    `tfsdk:"target_users"`
}

func (m *ProxyPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *ProxyPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *ProxyPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *ProxyPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *ProxyPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *ProxyPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *ProxyPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *ProxyPolicyModel) GetGroups() types.Set   { return m.Groups }

// SetProxyPolicyAttributes populates the TF schema data from a proxy policy
func SetProxyPolicyAttributes(ctx context.Context, schema *ProxyPolicyModel, apiPolicy *policies.ProxyPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)

	// See comment in SetBasePolicyAttributes that explains this conditional
	// logic
	if !schema.Environments.IsNull() || len(apiPolicy.GetEnvironments()) != 0 || modelIsDataSource {
		schema.Environments = policy.FlattenPolicyEnvironments(ctx, apiPolicy.GetEnvironments())
	}
	if !schema.Targets.IsNull() || len(apiPolicy.GetTargets()) != 0 || modelIsDataSource {
		schema.Targets = policy.FlattenPolicyTargets(ctx, apiPolicy.GetTargets())
	}
	if !schema.TargetUsers.IsNull() || len(apiPolicy.GetTargetUsers()) != 0 || modelIsDataSource {
		schema.TargetUsers = policy.FlattenPolicyTargetUsers(ctx, apiPolicy.GetTargetUsers())
	}
}

func ExpandProxyPolicy(ctx context.Context, schema *ProxyPolicyModel) *policies.ProxyPolicy {
	p := new(policies.ProxyPolicy)
	p.Name = schema.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, schema.Description)
	p.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, schema.Subjects))
	p.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, schema.Groups))
	p.Environments = bastionzero.PtrTo(policy.ExpandPolicyEnvironments(ctx, schema.Environments))
	p.Targets = bastionzero.PtrTo(policy.ExpandPolicyTargets(ctx, schema.Targets))
	p.TargetUsers = bastionzero.PtrTo(policy.ExpandPolicyTargetUsers(ctx, schema.TargetUsers))

	return p
}

func makeProxyPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.Proxy)
	attributes["environments"] = policy.PolicyEnvironmentsAttribute()
	attributes["targets"] = policy.PolicyTargetsAttribute([]targettype.TargetType{
		targettype.Db,
		targettype.Web,
	})
	attributes["target_users"] = schema.SetAttribute{
		Description: "Set of Database usernames that this policy applies to. " +
			"These usernames only affect policy decisions involving Db targets that have the SplitCert feature enabled.",
		ElementType: types.StringType,
		Optional:    true,
	}

	return attributes
}
