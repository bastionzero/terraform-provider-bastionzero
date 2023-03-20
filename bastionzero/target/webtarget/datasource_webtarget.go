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

func NewWebTargetDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[webTargetModel, targets.WebTarget]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[webTargetModel, targets.WebTarget]{
				RecordSchema: makeWebTargetDataSourceSchema(
					&target.BaseTargetDataSourceAttributeOptions{
						IsIDRequired:   true,
						IsNameComputed: true,
					}),
				ResultAttributeName: "web_target",
				PrettyAttributeName: "Web target",
				FlattenAPIModel: func(ctx context.Context, apiObject *targets.WebTarget) (state *webTargetModel, diags diag.Diagnostics) {
					state = new(webTargetModel)
					setWebTargetAttributes(ctx, state, apiObject)
					return
				},
				Description: "Get information about a specific Web target in your BastionZero organization.",
			},
			GetAPIModel: func(ctx context.Context, tfModel webTargetModel, client *bastionzero.Client) (*targets.WebTarget, error) {
				target, _, err := client.Targets.GetWebTarget(ctx, tfModel.ID.ValueString())
				return target, err
			},
		},
	)
}
