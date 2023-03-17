package environment

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/listdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewEnvironmentsDataSource() datasource.DataSource {
	return listdatasource.NewListDataSource(&listdatasource.ListDataSourceConfig[environmentModel, environments.Environment]{
		RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeEnvironmentResourceSchema(), nil),
		ResultAttributeName: "environments",
		PrettyAttributeName: "environments",
		FlattenAPIModel: func(ctx context.Context, apiObject environments.Environment) (envState *environmentModel, diags diag.Diagnostics) {
			envState = new(environmentModel)
			setEnvironmentAttributes(ctx, envState, &apiObject)
			return envState, diags
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]environments.Environment, error) {
			environments, _, err := client.Environments.ListEnvironments(ctx)
			return environments, err
		},
		Description: "Get a list of all environments in your BastionZero organization.",
	})
}
