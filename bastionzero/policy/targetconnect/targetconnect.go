package targetconnect

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TargetConnectPolicyModel maps the target connect policy schema data.
type TargetConnectPolicyModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	Subjects     types.Set    `tfsdk:"subjects"`
	Groups       types.Set    `tfsdk:"groups"`
	Environments types.Set    `tfsdk:"environments"`
	Targets      types.Set    `tfsdk:"targets"`
	TargetUsers  types.Set    `tfsdk:"target_users"`
	Verbs        types.Set    `tfsdk:"verbs"`
}

func (m *TargetConnectPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *TargetConnectPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *TargetConnectPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *TargetConnectPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *TargetConnectPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *TargetConnectPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *TargetConnectPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *TargetConnectPolicyModel) GetGroups() types.Set   { return m.Groups }

// SetTargetConnectPolicyAttributes populates the TF schema data from a target
// connect policy
func SetTargetConnectPolicyAttributes(ctx context.Context, schema *TargetConnectPolicyModel, apiPolicy *policies.TargetConnectPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)

	// See comment in SetBasePolicyAttributes that explains this conditional
	// logic
	if !schema.Environments.IsNull() || len(apiPolicy.GetEnvironments()) != 0 || modelIsDataSource {
		schema.Environments = policy.FlattenPolicyEnvironments(ctx, apiPolicy.GetEnvironments())
	}
	if !schema.Targets.IsNull() || len(apiPolicy.GetTargets()) != 0 || modelIsDataSource {
		schema.Targets = policy.FlattenPolicyTargets(ctx, apiPolicy.GetTargets())
	}

	// By def. of schema, these values cannot be null so just accept whatever
	// the refreshed value is
	if modelIsDataSource {
		schema.TargetUsers = policy.FlattenPolicyTargetUsers(ctx, apiPolicy.GetTargetUsers()[:len(apiPolicy.GetTargetUsers())-1])
	} else {
		schema.TargetUsers = policy.FlattenPolicyTargetUsers(ctx, apiPolicy.GetTargetUsers())
	}
	schema.Verbs = FlattenPolicyVerbs(ctx, apiPolicy.GetVerbs())

}

func ExpandTargetConnectPolicy(ctx context.Context, schema *TargetConnectPolicyModel) *policies.TargetConnectPolicy {
	p := new(policies.TargetConnectPolicy)
	p.Name = schema.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, schema.Description)
	p.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, schema.Subjects))
	p.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, schema.Groups))
	p.Environments = bastionzero.PtrTo(policy.ExpandPolicyEnvironments(ctx, schema.Environments))
	p.Targets = bastionzero.PtrTo(policy.ExpandPolicyTargets(ctx, schema.Targets))
	p.TargetUsers = bastionzero.PtrTo(policy.ExpandPolicyTargetUsers(ctx, schema.TargetUsers))
	p.Verbs = bastionzero.PtrTo(ExpandPolicyVerbs(ctx, schema.Verbs))

	return p
}

func ExpandPolicyVerbs(ctx context.Context, tfSet types.Set) []policies.Verb {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.Verb {
		return policies.Verb{Type: verbtype.VerbType(m)}
	})
}

func FlattenPolicyVerbs(ctx context.Context, apiObject []policies.Verb) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.Verb) attr.Value {
		return types.StringValue(string(m.Type))
	})
}

func makeTargetConnectPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.TargetConnect)
	attributes["environments"] = policy.PolicyEnvironmentsAttribute()
	attributes["targets"] = policy.PolicyTargetsAttribute([]targettype.TargetType{
		targettype.Bzero,
		targettype.DynamicAccessConfig,
	})
	attributes["target_users"] = schema.SetAttribute{
		Description: "Set of Unix usernames that this policy applies to.",
		ElementType: types.StringType,
		Required:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
	attributes["verbs"] = schema.SetAttribute{
		Required:    true,
		Description: fmt.Sprintf("Set of actions allowed by this policy %s.", internal.PrettyOneOf(verbtype.VerbTypeValues())),
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(stringvalidator.OneOf(bastionzero.ToStringSlice(verbtype.VerbTypeValues())...)),
			setvalidator.SizeAtLeast(1),
		},
	}

	return attributes
}
