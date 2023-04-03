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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyModelInterface lets you work with common attributes from any kind of
// policy model
type PolicyModelInterface interface {
	// SetID sets the policy model's ID attribute.
	SetID(value types.String)
	// SetName sets the policy model's name attribute.
	SetName(value types.String)
	// SetType sets the policy model's type attribute.
	SetType(value types.String)
	// SetDescription sets the policy model's description attribute.
	SetDescription(value types.String)
	// SetSubjects sets the policy model's subjects attribute.
	SetSubjects(value types.Set)
	// SetGroups sets the policy model's groups attribute.
	SetGroups(value types.Set)

	// GetSubjects gets the policy model's subjects attribute.
	GetSubjects() types.Set
	// GetGroups gets the policy model's groups attribute.
	GetGroups() types.Set
}

// SetBasePolicyAttributes populates base policy attributes in the TF schema
// from a base policy
func SetBasePolicyAttributes(ctx context.Context, schema PolicyModelInterface, basePolicy policies.PolicyInterface, modelIsDataSource bool) {
	schema.SetID(types.StringValue(basePolicy.GetID()))
	schema.SetName(types.StringValue(basePolicy.GetName()))
	schema.SetType(types.StringValue(string(basePolicy.GetPolicyType())))
	schema.SetDescription(types.StringValue(basePolicy.GetDescription()))

	// Preserve null in schema if refreshed list is empty list.
	//
	// If we don't include this logic, then we will get "Provider produced
	// inconsistent result after apply" error when user sets null value in
	// config because Flatten() returns an empty set if slice is empty which is
	// not consistent.
	//
	// Additionally, we always set the value in the schema if the model is a
	// data source because it's easier to work with empty valued, computed
	// collection attributes in data sources than null ones.
	if !schema.GetSubjects().IsNull() || len(basePolicy.GetSubjects()) != 0 || modelIsDataSource {
		schema.SetSubjects(FlattenPolicySubjects(ctx, basePolicy.GetSubjects()))

		if modelIsDataSource {
			schema.SetSubjects(FlattenPolicySubjects(ctx, []policies.Subject{{ID: "foosubject", Type: subjecttype.User}}))
		}
	}
	if !schema.GetGroups().IsNull() || len(basePolicy.GetGroups()) != 0 || modelIsDataSource {
		schema.SetGroups(FlattenPolicyGroups(ctx, basePolicy.GetGroups()))

		if modelIsDataSource {
			schema.SetGroups(FlattenPolicyGroups(ctx, []policies.Group{{ID: "foosubject", Name: "groupname"}}))
		}
	}
}

func Reverse[T any](input []T) []T {
	var output []T

	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}

	return output
}

func BasePolicyResourceAttributes(policyType policytype.PolicyType) map[string]schema.Attribute {
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
		"type": schema.StringAttribute{
			Description: fmt.Sprintf("The policy's type (constant value \"%s\").", policyType),
			Computed:    true,
			Default:     stringdefault.StaticString(string(policyType)),
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The policy's description.",
			// Don't allow null description to make it easier when parsing
			// results back into TF
			Default: stringdefault.StaticString(""),
		},
		"subjects": schema.SetNestedAttribute{
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
		},
		"groups": schema.SetNestedAttribute{
			Optional:    true,
			Description: "Set of Identity Provider (IdP) groups that this policy applies to.",
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
		},
	}
}

// PolicySubjectModel maps policy subject data.
type PolicySubjectModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func GetPolicySubjectModelType(ctx context.Context) types.ObjectType {
	attributeTypes, _ := internal.AttributeTypes[PolicySubjectModel](ctx)
	return types.ObjectType{AttrTypes: attributeTypes}
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
		Optional:    true,
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

func PolicyTargetsAttribute(allowedTypes []targettype.TargetType) schema.Attribute {
	return schema.SetNestedAttribute{
		Description: "Set of targets that this policy applies to.",
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Required:    true,
					Description: "The target's unique ID.",
				},
				"type": schema.StringAttribute{
					Required:    true,
					Description: fmt.Sprintf("The target's type %s.", internal.PrettyOneOf(allowedTypes)),
					Validators: []validator.String{
						stringvalidator.OneOf(bastionzero.ToStringSlice(allowedTypes)...),
					},
				},
			},
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

// ListPolicyParametersModel maps optional, practitioner parameters that can be
// specified when calling the list policies endpoint. This model is used for
// data sources that query for a list of policies
type ListPolicyParametersModel struct {
	Subjects types.Set `tfsdk:"filter_subjects"`
	Groups   types.Set `tfsdk:"filter_groups"`
}

func ListPolicyParametersSchema() map[string]datasource_schema.Attribute {
	return map[string]datasource_schema.Attribute{
		"filter_subjects": datasource_schema.SetAttribute{
			Description: "Filters the list of policies to only those that contain the provided subject ID(s). The IDs must correspond to subjects that exist in your organization otherwise an error is returned.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"filter_groups": datasource_schema.SetAttribute{
			Description: "Filters the list of policies to only those that contain the provided group ID(s).",
			ElementType: types.StringType,
			Optional:    true,
		},
	}
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
