package proxy

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// proxyPolicyModel maps the proxy policy schema data.
type proxyPolicyModel struct {
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

func (m *proxyPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *proxyPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *proxyPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *proxyPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *proxyPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *proxyPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *proxyPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *proxyPolicyModel) GetGroups() types.Set   { return m.Groups }

// setProxyPolicyAttributes populates the TF schema data from a proxy policy
func setProxyPolicyAttributes(ctx context.Context, schema *proxyPolicyModel, apiPolicy *policies.ProxyPolicy, modelIsDataSource bool) {
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
