package environment

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &environmentsDataSource{}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &environmentsDataSource{}
}

// environmentsDataSource is the data source implementation.
type environmentsDataSource struct {
	client *bastionzero.Client
}

// environmentsDataSourceModel describes the environments data source data model.
type environmentsDataSourceModel struct {
	Environments []environmentModel `tfsdk:"environments"`
}

// Metadata returns the users data source type name.
func (d *environmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

// Schema defines the schema for the environments data source.
func (d *environmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a list of all environments in your BastionZero organization.",
		Attributes: map[string]schema.Attribute{
			"environments": schema.ListNestedAttribute{
				Description: "List of environments in your organization.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: internal.ResourceSchemaToDataSourceSchema(makeEnvironmentResourceSchema(), nil),
				},
			},
		},
	}
}

// Configure adds the provider configured BastionZero API client to the data
// source.
func (d *environmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*bastionzero.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source configure type",
			fmt.Sprintf("Expected *bastionzero.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the users Terraform state with the latest data.
func (d *environmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state environmentsDataSourceModel

	// Query BastionZero for environments
	tflog.Debug(ctx, "Querying for environments")
	environments, _, err := d.client.Environments.ListEnvironments(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list environments",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Queried for environments", map[string]any{"num_environments": len(environments)})

	// Map response body to model
	for _, env := range environments {
		var envState environmentModel
		// TODO: Fix setEnvironmentAttribute to not throw error
		resp.Diagnostics.Append(setEnvironmentAttributes(ctx, &envState, &env)...)

		state.Environments = append(state.Environments, envState)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
