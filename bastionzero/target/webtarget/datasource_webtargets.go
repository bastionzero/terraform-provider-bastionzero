package webtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewWebTargetsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[webTargetModel, targets.WebTarget]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[webTargetModel, targets.WebTarget]{
			RecordSchema: makeWebTargetDataSourceSchema(
				&target.BaseTargetDataSourceAttributeOptions{
					IsIDComputed:   true,
					IsNameComputed: true,
				}),
			MetadataTypeName:    "web_targets",
			ResultAttributeName: "targets",
			PrettyAttributeName: "Web targets",
			FlattenAPIModel: func(ctx context.Context, apiObject *targets.WebTarget) (state *webTargetModel, diags diag.Diagnostics) {
				state = new(webTargetModel)
				setWebTargetAttributes(ctx, state, apiObject)
				return
			},
			MarkdownDescription: "Get a list of all Web targets in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.WebTarget, error) {
			targets, _, err := client.Targets.ListWebTargets(ctx)
			return targets, err
		},
	})
}
