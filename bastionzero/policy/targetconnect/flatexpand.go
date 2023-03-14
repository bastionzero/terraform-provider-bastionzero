package targetconnect

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExpandPolicyTargetUsers(ctx context.Context, tfSet types.Set) *[]policies.PolicyTargetUser {
	targetUsers := internal.ExpandFrameworkStringValueSet(ctx, tfSet)

	apiObject := make([]policies.PolicyTargetUser, len(targetUsers))
	for i, username := range targetUsers {
		apiObject[i] = policies.PolicyTargetUser{Username: username}
	}

	return &apiObject
}

func FlattenPolicyTargetUsers(ctx context.Context, apiObject *[]policies.PolicyTargetUser) types.Set {
	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(types.StringType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, v := range *apiObject {
		elements[i] = types.StringValue(v.Username)
	}

	return types.SetValueMust(types.StringType, elements)
}

func ExpandPolicyVerbs(ctx context.Context, tfSet types.Set) *[]policies.Verb {
	policyVerbs := internal.ExpandFrameworkStringValueSet(ctx, tfSet)

	apiObject := make([]policies.Verb, len(policyVerbs))
	for i, verbType := range policyVerbs {
		apiObject[i] = policies.Verb{Type: verbtype.VerbType(verbType)}
	}

	return &apiObject
}

func FlattenPolicyVerbs(ctx context.Context, apiObject *[]policies.Verb) types.Set {
	if apiObject == nil || len(*apiObject) == 0 {
		return types.SetNull(types.StringType)
	}

	elements := make([]attr.Value, len(*apiObject))
	for i, v := range *apiObject {
		elements[i] = types.StringValue(string(v.Type))
	}

	return types.SetValueMust(types.StringType, elements)
}
