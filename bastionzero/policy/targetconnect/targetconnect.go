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
