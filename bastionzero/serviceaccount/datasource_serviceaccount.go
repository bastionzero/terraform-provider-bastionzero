package serviceaccount

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewServiceAccountDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[serviceAccountModel, serviceaccounts.ServiceAccount]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[serviceAccountModel, serviceaccounts.ServiceAccount]{
				RecordSchema:        makeServiceAccountDataSourceSchema(true),
				MetadataTypeName:    "service_account",
				PrettyAttributeName: "service account",
				FlattenAPIModel: func(ctx context.Context, apiObject *serviceaccounts.ServiceAccount, state *serviceAccountModel) (diags diag.Diagnostics) {
					setServiceAccountAttributes(ctx, state, apiObject)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel serviceAccountModel, client *bastionzero.Client) (*serviceaccounts.ServiceAccount, error) {
					serviceAccount, _, err := client.ServiceAccounts.GetServiceAccount(ctx, tfModel.ID.ValueString())
					return serviceAccount, err
				},
				Description: "Get information on a service account in your BastionZero organization.",
			},
		},
	)
}
