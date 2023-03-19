package bzerotarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewBzeroTargetDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(&bzdatasource.SingleDataSourceConfig[bzeroTargetModel, targets.BzeroTarget]{
		RecordSchema:        makeBzeroTargetDataSourceSchema(true),
		ResultAttributeName: "bzero_target",
		PrettyAttributeName: "Bzero target",
		FlattenAPIModel: func(ctx context.Context, apiObject *targets.BzeroTarget) (state *bzeroTargetModel, diags diag.Diagnostics) {
			state = new(bzeroTargetModel)
			setBzeroTargetAttributes(ctx, state, apiObject)
			return
		},
		GetAPIModel: func(ctx context.Context, client *bastionzero.Client, id string) (*targets.BzeroTarget, error) {
			targets, _, err := client.Targets.GetBzeroTarget(ctx, id)
			return targets, err
		},
		Description: "Get information about a specific Bzero target in your BastionZero organization.",
	})

}
