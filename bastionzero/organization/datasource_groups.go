package organization

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &groupsDataSource{}

func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

// groupsDataSource is the data source implementation.
type groupsDataSource struct {
	client *bastionzero.Client
}

// groupsDataSourceModel describes the groups data source data model.
type groupsDataSourceModel struct {
	Groups []groupModel `tfsdk:"groups"`
}

// groupModel maps group schema data.
type groupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Metadata returns the groups data source type name.
func (d *groupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

// Schema defines the schema for the groups data source.
func (d *groupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a list of all groups in your BastionZero organization.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "List of groups in your organization.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The group's unique ID, as specified by the Identity Provider in which it is configured.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The group's name.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured BastionZero API client to the data
// source.
func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupsDataSourceModel

	// Query BastionZero for users
	tflog.Debug(ctx, "Querying for groups")
	groups, _, err := d.client.Organization.ListGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list groups",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Queried for groups", map[string]any{"num_groups": len(groups)})

	// Map response body to model
	for _, group := range groups {
		groupState := groupModel{
			ID:   types.StringValue(group.ID),
			Name: types.StringValue(group.Name),
		}
		state.Groups = append(state.Groups, groupState)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
