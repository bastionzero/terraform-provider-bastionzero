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

func NewTargetConnectPolicyDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[TargetConnectPolicyModel, policies.TargetConnectPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[TargetConnectPolicyModel, policies.TargetConnectPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeTargetConnectPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "targetconnect_policy",
				PrettyAttributeName: "target connect policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.TargetConnectPolicy, state *TargetConnectPolicyModel) (diags diag.Diagnostics) {
					SetTargetConnectPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel TargetConnectPolicyModel, client *bastionzero.Client) (*policies.TargetConnectPolicy, error) {
					policy, _, err := client.Policies.GetTargetConnectPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				Description: "Get information on a BastionZero target connect policy.",
			},
		},
	)
}
