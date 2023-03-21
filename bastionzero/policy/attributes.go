package policy

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicySubjectModel maps policy subject data.
type PolicySubjectModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func GetPolicySubjectModelType(ctx context.Context) types.ObjectType {
	attributeTypes, _ := internal.AttributeTypes[PolicySubjectModel](ctx)
	return types.ObjectType{AttrTypes: attributeTypes}
}

func PolicySubjectsAttribute(ctx context.Context) schema.Attribute {
	return schema.SetNestedAttribute{
		Computed:    true,
		Optional:    true,
		Description: "Set of subjects that this policy applies to.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Required:    true,
					Description: "The subject's unique ID.",
				},
				"type": schema.StringAttribute{
					Required:    true,
					Description: fmt.Sprintf("The subject's type %s.", internal.PrettyOneOf(subjecttype.SubjectTypeValues())),
					Validators: []validator.String{
						stringvalidator.OneOf(bastionzero.ToStringSlice(subjecttype.SubjectTypeValues())...),
					},
				},
			},
		},
		PlanModifiers: []planmodifier.Set{
			bzplanmodifier.SetDefaultValue(types.SetValueMust(GetPolicySubjectModelType(ctx), []attr.Value{})),
		},
	}
}

func ExpandPolicySubjects(ctx context.Context, tfSet types.Set) []policies.Subject {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m PolicySubjectModel) policies.Subject {
		return policies.Subject{
			ID:   m.ID.ValueString(),
			Type: subjecttype.SubjectType(m.Type.ValueString()),
		}
	})
}

func FlattenPolicySubjects(ctx context.Context, apiObject []policies.Subject) types.Set {
	elementType := GetPolicySubjectModelType(ctx)
	attributeTypes := elementType.AttrTypes
	return internal.FlattenFrameworkSet(ctx, elementType, apiObject, func(m policies.Subject) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(m.ID),
			"type": types.StringValue(string(m.Type)),
		})
	})
}

// PolicyGroupModel maps policy group data.
type PolicyGroupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func GetPolicyGroupModelType(ctx context.Context) types.ObjectType {
	attributeTypes, _ := internal.AttributeTypes[PolicyGroupModel](ctx)
	return types.ObjectType{AttrTypes: attributeTypes}
}

func PolicyGroupsAttribute(ctx context.Context) schema.Attribute {
	return schema.SetNestedAttribute{
		Computed:    true,
		Optional:    true,
		Description: "Set of IdP groups that this policy applies to.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Required:    true,
					Description: "The group's unique ID.",
				},
				"name": schema.StringAttribute{
					Required:    true,
					Description: "The group's name.",
				},
			},
		},
		PlanModifiers: []planmodifier.Set{
			bzplanmodifier.SetDefaultValue(types.SetValueMust(GetPolicyGroupModelType(ctx), []attr.Value{})),
		},
	}
}

func ExpandPolicyGroups(ctx context.Context, tfSet types.Set) []policies.Group {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m PolicyGroupModel) policies.Group {
		return policies.Group{
			ID:   m.ID.ValueString(),
			Name: m.Name.ValueString(),
		}
	})
}

func FlattenPolicyGroups(ctx context.Context, apiObject []policies.Group) types.Set {
	elementType := GetPolicyGroupModelType(ctx)
	attributeTypes := elementType.AttrTypes
	return internal.FlattenFrameworkSet(ctx, elementType, apiObject, func(m policies.Group) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(m.ID),
			"name": types.StringValue(m.Name),
		})
	})
}

func PolicyEnvironmentsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Description: "Set of environments that this policy applies to.",
		ElementType: types.StringType,
		Computed:    true,
		Optional:    true,
		PlanModifiers: []planmodifier.Set{
			bzplanmodifier.SetDefaultValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
	}
}

func ExpandPolicyEnvironments(ctx context.Context, tfSet types.Set) []policies.Environment {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m string) policies.Environment {
		return policies.Environment{ID: m}
	})
}

func FlattenPolicyEnvironments(ctx context.Context, apiObject []policies.Environment) types.Set {
	return internal.FlattenFrameworkSet(ctx, types.StringType, apiObject, func(m policies.Environment) attr.Value {
		return types.StringValue(m.ID)
	})
}

// PolicyTargetModel maps policy target data.
type PolicyTargetModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func GetPolicyTargetModelType(ctx context.Context) types.ObjectType {
	attributeTypes, _ := internal.AttributeTypes[PolicyTargetModel](ctx)
	return types.ObjectType{AttrTypes: attributeTypes}
}

func PolicyTargetsAttribute(ctx context.Context) schema.Attribute {
	return schema.SetNestedAttribute{
		Description: "Set of targets that this policy applies to.",
		Computed:    true,
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Required:    true,
					Description: "The target's unique ID.",
				},
				"type": schema.StringAttribute{
					Required:    true,
					Description: fmt.Sprintf("The target's type %s.", internal.PrettyOneOf(targettype.TargetTypeValues())),
					Validators: []validator.String{
						stringvalidator.OneOf(bastionzero.ToStringSlice(targettype.TargetTypeValues())...),
					},
				},
			},
		},
		PlanModifiers: []planmodifier.Set{
			bzplanmodifier.SetDefaultValue(types.SetValueMust(GetPolicyTargetModelType(ctx), []attr.Value{})),
		},
	}
}

func ExpandPolicyTargets(ctx context.Context, tfSet types.Set) []policies.Target {
	return internal.ExpandFrameworkSet(ctx, tfSet, func(m PolicyTargetModel) policies.Target {
		return policies.Target{
			ID:   m.ID.ValueString(),
			Type: targettype.TargetType(m.Type.ValueString()),
		}
	})
}

func FlattenPolicyTargets(ctx context.Context, apiObject []policies.Target) types.Set {
	elementType := GetPolicyTargetModelType(ctx)
	attributeTypes := elementType.AttrTypes
	return internal.FlattenFrameworkSet(ctx, elementType, apiObject, func(m policies.Target) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(m.ID),
			"type": types.StringValue(string(m.Type)),
		})
	})
}

func PolicyTypeAttribute(policyType policytype.PolicyType) schema.Attribute {
	return schema.StringAttribute{
		Description: fmt.Sprintf("The policy's type (constant value \"%s\").", policyType),
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			bzplanmodifier.StringDefaultValue(types.StringValue(string(policyType))),
		},
	}
}
