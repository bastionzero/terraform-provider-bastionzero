package environment

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewEnvironmentDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[environmentModel, environments.Environment]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[environmentModel, environments.Environment]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeEnvironmentResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "environment",
				PrettyAttributeName: "environment",
				FlattenAPIModel: func(ctx context.Context, apiObject *environments.Environment, state *environmentModel) (diags diag.Diagnostics) {
					setEnvironmentAttributes(ctx, state, apiObject)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel environmentModel, client *bastionzero.Client) (*environments.Environment, error) {
					env, _, err := client.Environments.GetEnvironment(ctx, tfModel.ID.ValueString())
					return env, err
				},
				Description: "Get information on a BastionZero environment.",
			},
		},
	)
}
