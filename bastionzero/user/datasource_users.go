package user

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &usersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *bastionzero.Client
}

// usersDataSourceModel describes the users data source data model.
type usersDataSourceModel struct {
	Users []userModel `tfsdk:"users"`
}

// userModel maps user schema data.
type userModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	FullName       types.String `tfsdk:"full_name"`
	Email          types.String `tfsdk:"email"`
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
	TimeCreated    types.String `tfsdk:"time_created"`
	LastLogin      types.String `tfsdk:"last_login"`
}

// Metadata returns the users data source type name.
func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the users data source.
func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a list of all users in your BastionZero organization.",
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Description: "List of users in your organization.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The user's unique ID.",
						},
						"organization_id": schema.StringAttribute{
							Computed:    true,
							Description: "The user's organization's ID.",
						},
						"full_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's full name.",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "The user's email address.",
						},
						"is_admin": schema.BoolAttribute{
							Computed:    true,
							Description: "If true, the user is an administrator. False otherwise.",
						},
						"time_created": schema.StringAttribute{
							Computed:    true,
							Description: "The time this user was created in BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
						},
						"last_login": schema.StringAttribute{
							Computed:    true,
							Description: "The time this user last logged into BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured BastionZero API client to the data
// source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	// Query BastionZero for users
	tflog.Debug(ctx, "Querying for users")
	users, _, err := d.client.Users.ListUsers(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list users",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Queried for users", map[string]any{"num_users": len(users)})

	// Map response body to model
	for _, user := range users {
		userState := userModel{
			ID:             types.StringValue(user.ID),
			OrganizationID: types.StringValue(user.OrganizationID),
			FullName:       types.StringValue(user.FullName),
			Email:          types.StringValue(user.Email),
			IsAdmin:        types.BoolValue(user.IsAdmin),
			TimeCreated:    types.StringValue(user.TimeCreated.UTC().Format(time.RFC3339)),
		}

		if user.LastLogin != nil {
			userState.LastLogin = types.StringValue(user.LastLogin.UTC().Format(time.RFC3339))
		} else {
			userState.LastLogin = types.StringNull()
		}

		state.Users = append(state.Users, userState)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
