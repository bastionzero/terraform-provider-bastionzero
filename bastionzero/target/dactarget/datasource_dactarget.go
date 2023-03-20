package dactarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewDacTargetDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[dacTargetModel, targets.DynamicAccessConfiguration]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[dacTargetModel, targets.DynamicAccessConfiguration]{
				RecordSchema: makeDacTargetDataSourceSchema(
					&dacTargetDataSourceAttributeOptions{
						IsIDRequired: true,
					}),
				ResultAttributeName: "dac_target",
				PrettyAttributeName: "Dynamic access configuration target",
				FlattenAPIModel: func(ctx context.Context, apiObject *targets.DynamicAccessConfiguration) (state *dacTargetModel, diags diag.Diagnostics) {
					state = new(dacTargetModel)
					setDacTargetAttributes(ctx, state, apiObject)
					return
				},
				Description: "Get information about a specific dynamic access configuration (DAC) target in your BastionZero organization.",
			},
			GetAPIModel: func(ctx context.Context, tfModel dacTargetModel, client *bastionzero.Client) (*targets.DynamicAccessConfiguration, error) {
				target, _, err := client.Targets.GetDynamicAccessConfiguration(ctx, tfModel.ID.ValueString())
				return target, err
			},
		},
	)
}
