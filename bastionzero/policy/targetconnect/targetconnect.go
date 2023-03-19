package targetconnect

import (
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makeTargetConnectPolicyResourceSchema() map[string]schema.Attribute {
	res := map[string]schema.Attribute{
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
			Description: "The policy's description.",
		},
		"subjects":     policy.PolicySubjectsAttribute(),
		"groups":       policy.PolicyGroupsAttribute(),
		"environments": policy.PolicyEnvironmentsAttribute(),
		"targets":      policy.PolicyTargetsAttribute(),
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

	return res
}
