package dactarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDacTargetsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[dacTargetModel, targets.DynamicAccessConfiguration]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[dacTargetModel, targets.DynamicAccessConfiguration]{
			RecordSchema: makeDacTargetDataSourceSchema(
				&dacTargetDataSourceAttributeOptions{
					IsIDComputed: true,
				}),
			ResultAttributeName: "dac_targets",
			PrettyAttributeName: "dynamic access configuration targets",
			FlattenAPIModel: func(ctx context.Context, apiObject *targets.DynamicAccessConfiguration) (state *dacTargetModel, diags diag.Diagnostics) {
				state = new(dacTargetModel)
				setDacTargetAttributes(ctx, state, apiObject)
				return
			},
			Description: "Get a list of all dynamic access configuration (DAC) targets in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.DynamicAccessConfiguration, error) {
			targets, _, err := client.Targets.ListDynamicAccessConfigurations(ctx)
			return targets, err
		},
	})
}
