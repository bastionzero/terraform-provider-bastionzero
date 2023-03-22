package targetconnect

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// targetConnectPolicyModel maps the target connect policy schema data.
type targetConnectPolicyModel struct {
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

// setTargetConnectPolicyAttributes populates the TF schema data from a target
// connect policy
func setTargetConnectPolicyAttributes(ctx context.Context, schema *targetConnectPolicyModel, apiPolicy *policies.TargetConnectPolicy, modelIsDataSource bool) {
	schema.ID = types.StringValue(apiPolicy.ID)
	schema.Name = types.StringValue(apiPolicy.Name)
	schema.Type = types.StringValue(string(apiPolicy.GetPolicyType()))
	schema.Description = types.StringValue(apiPolicy.GetDescription())

	// Preserve null in schema if refreshed list is empty list.
	//
	// If we don't include this logic, then we will get "Provider produced
	// inconsistent result after apply" error when user sets null value in
	// config because Flatten() returns an empty set if slice is empty which is
	// not consistent.
	//
	// Additionally, we always set the value in the schema if the model is a
	// data source because it's easier to work with empty valued, computed
	// collection attributes in data sources then null.
	if !schema.Subjects.IsNull() || len(apiPolicy.GetSubjects()) != 0 || modelIsDataSource {
		schema.Subjects = policy.FlattenPolicySubjects(ctx, apiPolicy.GetSubjects())
	}
	if !schema.Groups.IsNull() || len(apiPolicy.GetGroups()) != 0 || modelIsDataSource {
		schema.Groups = policy.FlattenPolicyGroups(ctx, apiPolicy.GetGroups())
	}
	if !schema.Environments.IsNull() || len(apiPolicy.GetEnvironments()) != 0 || modelIsDataSource {
		schema.Environments = policy.FlattenPolicyEnvironments(ctx, apiPolicy.GetEnvironments())
	}
	if !schema.Targets.IsNull() || len(apiPolicy.GetTargets()) != 0 || modelIsDataSource {
		schema.Targets = policy.FlattenPolicyTargets(ctx, apiPolicy.GetTargets())
	}

	// By def. of schema, these values cannot be null so just accept whatever
	// the refreshed value is
	schema.TargetUsers = FlattenPolicyTargetUsers(ctx, apiPolicy.GetTargetUsers())
	schema.Verbs = FlattenPolicyVerbs(ctx, apiPolicy.GetVerbs())
}

func ExpandPolicyTargetUsers(ctx context.Context, tfSet types.Set) []policies.TargetUser {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.TargetUser {
		return policies.TargetUser{Username: m}
	})
}

func FlattenPolicyTargetUsers(ctx context.Context, apiObject []policies.TargetUser) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.TargetUser) attr.Value {
		return types.StringValue(m.Username)
	})
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

func makeTargetConnectPolicyResourceSchema(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The policy's unique ID.",
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "The policy's name.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"type": policy.PolicyTypeAttribute(policytype.TargetConnect),
		"description": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The policy's description.",
			PlanModifiers: []planmodifier.String{
				// Don't allow null description to make it easier when parsing
				// results back into TF
				bzplanmodifier.StringDefaultValue(types.StringValue("")),
			},
		},
		"subjects":     policy.PolicySubjectsAttribute(ctx),
		"groups":       policy.PolicyGroupsAttribute(ctx),
		"environments": policy.PolicyEnvironmentsAttribute(),
		"targets":      policy.PolicyTargetsAttribute(ctx),
		"target_users": schema.SetAttribute{
			Description: "Set of Unix usernames that this policy applies to.",
			ElementType: types.StringType,
			Required:    true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"verbs": schema.SetAttribute{
			Required:    true,
			Description: fmt.Sprintf("Set of actions allowed by this policy %s.", internal.PrettyOneOf(verbtype.VerbTypeValues())),
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(stringvalidator.OneOf(bastionzero.ToStringSlice(verbtype.VerbTypeValues())...)),
				setvalidator.SizeAtLeast(1),
			},
		},
	}
}
