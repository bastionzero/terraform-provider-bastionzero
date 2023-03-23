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
		&bzdatasource.SingleDataSourceConfig[targetConnectPolicyModel, policies.TargetConnectPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[targetConnectPolicyModel, policies.TargetConnectPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeTargetConnectPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "targetconnect_policy",
				PrettyAttributeName: "target connect policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.TargetConnectPolicy, state *targetConnectPolicyModel) (diags diag.Diagnostics) {
					setTargetConnectPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel targetConnectPolicyModel, client *bastionzero.Client) (*policies.TargetConnectPolicy, error) {
					env, _, err := client.Policies.GetTargetConnectPolicy(ctx, tfModel.ID.ValueString())
					return env, err
				},
				Description: "Get information on a BastionZero target connect policy.",
			},
		},
	)
}
