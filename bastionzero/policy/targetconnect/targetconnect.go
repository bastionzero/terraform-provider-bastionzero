package targetconnect

import (
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
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
		},
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
			Optional:    true,
		},
		"verbs": schema.SetAttribute{
			Optional:    true,
			Description: fmt.Sprintf("Set of actions allowed by this policy %s.", internal.PrettyOneOf(verbtype.VerbTypeValues())),
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(stringvalidator.OneOfCaseInsensitive(bastionzero.ToStringSlice(verbtype.VerbTypeValues())...)),
			},
		},
	}

	return res
}
