package jit

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// jitPolicyModel maps the JIT policy schema data.
type jitPolicyModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	Description           types.String `tfsdk:"description"`
	Subjects              types.Set    `tfsdk:"subjects"`
	Groups                types.Set    `tfsdk:"groups"`
	ChildPolicies         types.Set    `tfsdk:"child_policies"`
	AutomaticallyApproved types.Bool   `tfsdk:"auto_approved"`
	Duration              types.Int64  `tfsdk:"duration"`
}

func (m *jitPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *jitPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *jitPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *jitPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *jitPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *jitPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *jitPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *jitPolicyModel) GetGroups() types.Set   { return m.Groups }

// setJITPolicyAttributes populates the TF schema data from a JIT policy
func setJITPolicyAttributes(ctx context.Context, schema *jitPolicyModel, apiPolicy *policies.JITPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)
	// By def. of schema, ChildPolicies set cannot be null so just accept
	// whatever the refreshed value is
	schema.ChildPolicies = FlattenChildPolicies(ctx, apiPolicy.GetChildPolicies())
	schema.AutomaticallyApproved = types.BoolValue(apiPolicy.GetAutomaticallyApproved())
	schema.Duration = types.Int64Value(int64(apiPolicy.GetDuration()))
}

// ChildPolicyModel maps child policy data.
type ChildPolicyModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
}

func GetChildPolicyModelType(ctx context.Context) types.ObjectType {
	attributeTypes, _ := internal.AttributeTypes[ChildPolicyModel](ctx)
	return types.ObjectType{AttrTypes: attributeTypes}
}

func ExpandChildPolicies(ctx context.Context, tfSet types.Set) []string {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m ChildPolicyModel) string {
		return m.ID.ValueString()
	})
}

func FlattenChildPolicies(ctx context.Context, apiObject []policies.ChildPolicy) types.Set {
	elementType := GetChildPolicyModelType(ctx)
	attributeTypes := elementType.AttrTypes
	return internal.FlattenFrameworkSet(ctx, elementType, apiObject, func(m policies.ChildPolicy) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(m.ID),
			"type": types.StringValue(string(m.Type)),
			"name": types.StringValue(string(m.Name)),
		})
	})
}

func allowedChildPolicyTypes() []policytype.PolicyType {
	return []policytype.PolicyType{
		policytype.TargetConnect,
		policytype.Kubernetes,
		policytype.Proxy,
	}
}

func makeJITPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.JustInTime)
	attributes["child_policies"] = schema.SetNestedAttribute{
		Required:    true,
		Description: "Set of policies that a JIT policy applies to.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Required:    true,
					Description: "The policy's unique ID.",
				},
				"type": schema.StringAttribute{
					Computed:    true,
					Description: fmt.Sprintf("The policy's type %s.", internal.PrettyOneOf(allowedChildPolicyTypes())),
					Validators: []validator.String{
						stringvalidator.OneOf(bastionzero.ToStringSlice(allowedChildPolicyTypes())...),
					},
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "The policy's name.",
				},
			},
		},
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
	attributes["auto_approved"] = schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Description: "If true, then the policies created by this JIT policy will be automatically approved. " +
			"If false, then policies will only be created based on request and approval from reviewers (Defaults to false).",
		PlanModifiers: []planmodifier.Bool{
			bzplanmodifier.BoolDefaultValue(types.BoolValue(false)),
		},
	}
	attributes["duration"] = schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "The amount of time (in minutes) after which the access granted by this JIT policy will expire (Defaults to 1 hour).",
		PlanModifiers: []planmodifier.Int64{
			// Same default as the webapp
			bzplanmodifier.Int64DefaultValue(types.Int64Value(60)),
		},
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	}

	return attributes
}
