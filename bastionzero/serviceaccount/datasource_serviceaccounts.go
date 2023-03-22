package serviceaccount

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewServiceAccountsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[serviceAccountModel, serviceaccounts.ServiceAccount]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[serviceAccountModel, serviceaccounts.ServiceAccount]{
			RecordSchema:        makeServiceAccountDataSourceSchema(false),
			ResultAttributeName: "service_accounts",
			PrettyAttributeName: "service accounts",
			FlattenAPIModel: func(ctx context.Context, apiObject *serviceaccounts.ServiceAccount) (state *serviceAccountModel, diags diag.Diagnostics) {
				state = new(serviceAccountModel)
				setServiceAccountAttributes(ctx, state, apiObject)
				return
			},
			Description: "Get a list of all service accounts in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]serviceaccounts.ServiceAccount, error) {
			serviceAccounts, _, err := client.ServiceAccounts.ListServiceAccounts(ctx)
			return serviceAccounts, err
		},
	})
}
