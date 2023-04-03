package proxy

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewProxyPolicyDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[ProxyPolicyModel, policies.ProxyPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[ProxyPolicyModel, policies.ProxyPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeProxyPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "proxy_policy",
				PrettyAttributeName: "proxy policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.ProxyPolicy, state *ProxyPolicyModel) (diags diag.Diagnostics) {
					SetProxyPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel ProxyPolicyModel, client *bastionzero.Client) (*policies.ProxyPolicy, error) {
					policy, _, err := client.Policies.GetProxyPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				Description: "Get information on a BastionZero proxy policy.",
			},
		},
	)
}
