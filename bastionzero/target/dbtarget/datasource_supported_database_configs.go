package dbtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSupportedDatabaseConfigsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[types.Object, dbauthconfig.DatabaseAuthenticationConfig]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[types.Object, dbauthconfig.DatabaseAuthenticationConfig]{
			RecordSchema:        DatabaseAuthenticationConfigAttributes(),
			MetadataTypeName:    "supported_database_configs",
			ResultAttributeName: "configs",
			PrettyAttributeName: "supported database authentication configs",
			FlattenAPIModel: func(ctx context.Context, apiObject *dbauthconfig.DatabaseAuthenticationConfig) (state *types.Object, diags diag.Diagnostics) {
				return bastionzero.PtrTo(FlattenDatabaseAuthenticationConfig(ctx, apiObject)), diags
			},
			MarkdownDescription: "Get a list of all supported database authentication configs.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]dbauthconfig.DatabaseAuthenticationConfig, error) {
			configs, _, err := client.Targets.ListDatabaseAuthenticationConfigs(ctx)
			return configs, err
		},
	})
}
