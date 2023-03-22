package targetconnect

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewTargetConnectPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[targetConnectPolicyModel, policies.TargetConnectPolicy]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[targetConnectPolicyModel, policies.TargetConnectPolicy]{
			RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeTargetConnectPolicyResourceSchema(context.TODO()), nil),
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
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]policies.TargetConnectPolicy, error) {
			policies, _, err := client.Policies.ListTargetConnectPolicies(ctx, &policies.ListPolicyOptions{})
			return policies, err
		},
	})
}
