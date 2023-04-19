package bzerotarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewBzeroTargetsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[bzeroTargetModel, targets.BzeroTarget]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[bzeroTargetModel, targets.BzeroTarget]{
			RecordSchema: makeBzeroTargetDataSourceSchema(
				&target.BaseTargetDataSourceAttributeOptions{
					IsIDComputed:   true,
					IsNameComputed: true,
				}),
			MetadataTypeName:    "bzero_targets",
			ResultAttributeName: "targets",
			PrettyAttributeName: "Bzero targets",
			FlattenAPIModel: func(ctx context.Context, apiObject *targets.BzeroTarget) (state *bzeroTargetModel, diags diag.Diagnostics) {
				state = new(bzeroTargetModel)
				setBzeroTargetAttributes(ctx, state, apiObject)
				return
			},
			MarkdownDescription: "Get a list of all Bzero targets in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.BzeroTarget, error) {
			targets, _, err := client.Targets.ListBzeroTargets(ctx)
			return targets, err
		},
	})
}
