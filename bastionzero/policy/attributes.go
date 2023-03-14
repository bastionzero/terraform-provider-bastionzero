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

// policySubjectModel maps policy subject data.
type policySubjectModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func policySubjectsAttribute() schema.Attribute {
	return schema.SetNestedAttribute{
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

func expandPolicySubjects(ctx context.Context, tfList types.List) *[]policies.PolicySubject {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil
	}

	var data []policySubjectModel

	if diags := tfList.ElementsAs(ctx, &data, false); diags.HasError() {
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

func flattenPolicySubjects(ctx context.Context, apiObject *[]policies.PolicySubject) types.Set {
	attributeTypes, _ := internal.AttributeTypes[policySubjectModel](ctx)
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

// policyGroupModel maps policy group data.
type policyGroupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func PolicyGroupsAttribute() schema.Attribute {
	return schema.SetNestedAttribute{
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

func expandPolicyGroups(ctx context.Context, tfList types.List) *[]policies.PolicyGroup {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil
	}

	var data []policyGroupModel

	if diags := tfList.ElementsAs(ctx, &data, false); diags.HasError() {
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

func flattenPolicyGroups(ctx context.Context, apiObject *[]policies.PolicyGroup) types.Set {
	attributeTypes, _ := internal.AttributeTypes[policyGroupModel](ctx)
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

func policyEnvironmentsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Description: "Set of environments that this policy applies to.",
		ElementType: types.StringType,
	}
}

// policyTargetModel maps policy target data.
type policyTargetModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func policyTargetsAttribute() schema.Attribute {
	return schema.SetNestedAttribute{
		Description: "Set of targets that this policy applies to.",
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

func expandPolicyTargets(ctx context.Context, tfList types.List) *[]policies.PolicyTarget {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil
	}

	var data []policyTargetModel

	if diags := tfList.ElementsAs(ctx, &data, false); diags.HasError() {
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

func flattenPolicyTargets(ctx context.Context, apiObject *[]policies.PolicyTarget) types.Set {
	attributeTypes, _ := internal.AttributeTypes[policyTargetModel](ctx)
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
