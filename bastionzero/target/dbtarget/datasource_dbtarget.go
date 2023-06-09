package dbtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDbTargetDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[dbTargetModel, targets.DatabaseTarget]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[dbTargetModel, targets.DatabaseTarget]{
				RecordSchema: makeDbTargetDataSourceSchema(
					&target.BaseTargetDataSourceAttributeOptions{
						IsIDRequired:   true,
						IsNameComputed: true,
					}),
				MetadataTypeName:    "db_target",
				PrettyAttributeName: "Db target",
				FlattenAPIModel: func(ctx context.Context, apiObject *targets.DatabaseTarget, state *dbTargetModel) (diags diag.Diagnostics) {
					setDbTargetAttributes(ctx, state, apiObject)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel dbTargetModel, client *bastionzero.Client) (*targets.DatabaseTarget, error) {
					target, _, err := client.Targets.GetDatabaseTarget(ctx, tfModel.ID.ValueString())
					return target, err
				},
				MarkdownDescription: "Get information about a specific Db target in your BastionZero organization.",
			},
		},
	)
}
