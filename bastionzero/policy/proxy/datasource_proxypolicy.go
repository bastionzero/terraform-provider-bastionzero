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
		&bzdatasource.SingleDataSourceConfig[proxyPolicyModel, policies.ProxyPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[proxyPolicyModel, policies.ProxyPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeProxyPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "proxy_policy",
				PrettyAttributeName: "proxy policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.ProxyPolicy, state *proxyPolicyModel) (diags diag.Diagnostics) {
					setProxyPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel proxyPolicyModel, client *bastionzero.Client) (*policies.ProxyPolicy, error) {
					policy, _, err := client.Policies.GetProxyPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				Description: "Get information on a BastionZero proxy policy.",
			},
		},
	)
}
