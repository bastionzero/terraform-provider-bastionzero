package targetconnect

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

func NewTargetConnectPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSourceWithPractitionerParameters(
		&bzdatasource.ListDataSourceWithPractitionerParametersConfig[targetConnectPolicyModel, policy.ListPolicyParametersModel, policies.TargetConnectPolicy]{
			BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[targetConnectPolicyModel, policies.TargetConnectPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeTargetConnectPolicyResourceSchema(), nil),
				MetadataTypeName:    "targetconnect_policies",
				ResultAttributeName: "policies",
				PrettyAttributeName: "target connect policies",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.TargetConnectPolicy) (state *targetConnectPolicyModel, diags diag.Diagnostics) {
					state = new(targetConnectPolicyModel)
					setTargetConnectPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				Description: "Get a list of all target connect policies in your BastionZero organization.",
			},
			PractitionerParamsRecordSchema: policy.ListPolicyParametersSchema(),
			ListAPIModels: func(ctx context.Context, listParameters policy.ListPolicyParametersModel, client *bastionzero.Client) ([]policies.TargetConnectPolicy, error) {
				subjectsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Subjects), ",")
				groupsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Groups), ",")

				policies, _, err := client.Policies.ListTargetConnectPolicies(ctx, &policies.ListPolicyOptions{Subjects: subjectsFilter, Groups: groupsFilter})
				return policies, err
			},
		})
}
