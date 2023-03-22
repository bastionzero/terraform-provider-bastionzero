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

func NewEnvironmentsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[environmentModel, environments.Environment]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[environmentModel, environments.Environment]{
			RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeEnvironmentResourceSchema(), nil),
			ResultAttributeName: "environments",
			PrettyAttributeName: "environments",
			FlattenAPIModel: func(ctx context.Context, apiObject *environments.Environment) (state *environmentModel, diags diag.Diagnostics) {
				state = new(environmentModel)
				setEnvironmentAttributes(ctx, state, apiObject)
				return
			},
			Description: "Get a list of all environments in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]environments.Environment, error) {
			environments, _, err := client.Environments.ListEnvironments(ctx)
			return environments, err
		},
	})
}
