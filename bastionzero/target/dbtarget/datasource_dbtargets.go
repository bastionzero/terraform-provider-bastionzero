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

func NewDbTargetsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[dbTargetModel, targets.DatabaseTarget]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[dbTargetModel, targets.DatabaseTarget]{
			RecordSchema: makeDbTargetDataSourceSchema(
				&target.BaseTargetDataSourceAttributeOptions{
					IsIDComputed:   true,
					IsNameComputed: true,
				}),
			ResultAttributeName: "db_targets",
			PrettyAttributeName: "Db targets",
			FlattenAPIModel: func(ctx context.Context, apiObject *targets.DatabaseTarget) (state *dbTargetModel, diags diag.Diagnostics) {
				state = new(dbTargetModel)
				setDbTargetAttributes(ctx, state, apiObject)
				return
			},
			Description: "Get a list of all Db targets in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.DatabaseTarget, error) {
			targets, _, err := client.Targets.ListDatabaseTargets(ctx)
			return targets, err
		},
	})
}
