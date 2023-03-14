package policy

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicySubjectModel maps policy subject data.
type PolicySubjectModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func PolicySubjectsAttribute() schema.Attribute {
	return schema.SetNestedAttribute{
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
						stringvalidator.OneOfCaseInsensitive(bastionzero.ToStringSlice(subjecttype.SubjectTypeValues())...),
					},
				},
			},
		},
	}
}

func ExpandPolicySubjects(ctx context.Context, tfSet types.Set) *[]policies.PolicySubject {
	if tfSet.IsNull() || tfSet.IsUnknown() {
		return nil
	}

	var data []PolicySubjectModel

	if diags := tfSet.ElementsAs(ctx, &data, false); diags.HasError() {
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	apiObject := make([]policies.PolicySubject, len(data))
	for i, obj := range data {
		apiObject[i] = policies.PolicySubject{ID: obj.ID.ValueString(), Type: subjecttype.SubjectType(obj.Type.ValueString())}
	}

	return &apiObject
}

func FlattenPolicySubjects(ctx context.Context, apiObject *[]policies.PolicySubject) types.Set {
	attributeTypes, _ := internal.AttributeTypes[PolicySubjectModel](ctx)
	elementType := types.ObjectType{AttrTypes: attributeTypes}

	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(elementType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, obj := range *apiObject {
		elements[i] = types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(obj.ID),
			"type": types.StringValue(string(obj.Type)),
		})
	}

	return types.SetValueMust(elementType, elements)
}

// PolicyGroupModel maps policy group data.
type PolicyGroupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func PolicyGroupsAttribute() schema.Attribute {
	return schema.SetNestedAttribute{
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
	}
}

func ExpandPolicyGroups(ctx context.Context, tfSet types.Set) *[]policies.PolicyGroup {
	if tfSet.IsNull() || tfSet.IsUnknown() {
		return nil
	}

	var data []PolicyGroupModel

	if diags := tfSet.ElementsAs(ctx, &data, false); diags.HasError() {
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	apiObject := make([]policies.PolicyGroup, len(data))
	for i, obj := range data {
		apiObject[i] = policies.PolicyGroup{ID: obj.ID.ValueString(), Name: obj.Name.ValueString()}
	}

	return &apiObject
}

func FlattenPolicyGroups(ctx context.Context, apiObject *[]policies.PolicyGroup) types.Set {
	attributeTypes, _ := internal.AttributeTypes[PolicyGroupModel](ctx)
	elementType := types.ObjectType{AttrTypes: attributeTypes}

	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(elementType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, obj := range *apiObject {
		elements[i] = types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(obj.ID),
			"name": types.StringValue(obj.Name),
		})
	}

	return types.SetValueMust(elementType, elements)
}

func PolicyEnvironmentsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Description: "Set of environments that this policy applies to.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

func ExpandPolicyEnvironments(ctx context.Context, tfSet types.Set) *[]policies.PolicyEnvironment {
	envIds := internal.ExpandFrameworkStringValueSet(ctx, tfSet)

	apiObject := make([]policies.PolicyEnvironment, len(envIds))
	for i, id := range envIds {
		apiObject[i] = policies.PolicyEnvironment{ID: id}
	}

	return &apiObject
}

func FlattenPolicyEnvironments(ctx context.Context, apiObject *[]policies.PolicyEnvironment) types.Set {
	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(types.StringType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, v := range *apiObject {
		elements[i] = types.StringValue(v.ID)
	}

	return types.SetValueMust(types.StringType, elements)
}

// policyTargetModel maps policy target data.
type PolicyTargetModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func PolicyTargetsAttribute() schema.Attribute {
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
					Description: fmt.Sprintf("The target's type %s.", internal.PrettyOneOf(targettype.TargetTypeValues())),
					Validators: []validator.String{
						stringvalidator.OneOfCaseInsensitive(bastionzero.ToStringSlice(targettype.TargetTypeValues())...),
					},
				},
			},
		},
	}
}

func ExpandPolicyTargets(ctx context.Context, tfSet types.Set) *[]policies.PolicyTarget {
	if tfSet.IsNull() || tfSet.IsUnknown() {
		return nil
	}

	var data []PolicyTargetModel

	if diags := tfSet.ElementsAs(ctx, &data, false); diags.HasError() {
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	apiObject := make([]policies.PolicyTarget, len(data))
	for i, obj := range data {
		apiObject[i] = policies.PolicyTarget{ID: obj.ID.ValueString(), Type: targettype.TargetType(obj.Type.ValueString())}
	}

	return &apiObject
}

func FlattenPolicyTargets(ctx context.Context, apiObject *[]policies.PolicyTarget) types.Set {
	attributeTypes, _ := internal.AttributeTypes[PolicyTargetModel](ctx)
	elementType := types.ObjectType{AttrTypes: attributeTypes}

	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(elementType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, obj := range *apiObject {
		elements[i] = types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":   types.StringValue(obj.ID),
			"type": types.StringValue(string(obj.Type)),
		})
	}

	return types.SetValueMust(elementType, elements)
}
