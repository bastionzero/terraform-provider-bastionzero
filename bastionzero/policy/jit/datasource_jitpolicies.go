package jit

import (
	"context"
	"strings"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewJITPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSourceWithPractitionerParameters(
		&bzdatasource.ListDataSourceWithPractitionerParametersConfig[JITPolicyModel, policy.ListPolicyParametersModel, policies.JITPolicy]{
			BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[JITPolicyModel, policies.JITPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeJITPolicyResourceSchema(), nil),
				MetadataTypeName:    "jit_policies",
				ResultAttributeName: "policies",
				PrettyAttributeName: "JIT policies",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.JITPolicy) (state *JITPolicyModel, diags diag.Diagnostics) {
					state = new(JITPolicyModel)
					SetJITPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				Description: "Get a list of all JIT policies in your BastionZero organization.",
			},
			PractitionerParamsRecordSchema: policy.ListPolicyParametersSchema(),
			ListAPIModels: func(ctx context.Context, listParameters policy.ListPolicyParametersModel, client *bastionzero.Client) ([]policies.JITPolicy, error) {
				subjectsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Subjects), ",")
				groupsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Groups), ",")

				policies, _, err := client.Policies.ListJITPolicies(ctx, &policies.ListPolicyOptions{Subjects: subjectsFilter, Groups: groupsFilter})
				return policies, err
			},
		})
}
