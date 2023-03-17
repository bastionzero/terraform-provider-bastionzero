package user

import (
	"context"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/listdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func NewUsersDataSource() datasource.DataSource {
	return listdatasource.NewListDataSource(&listdatasource.ListDataSourceConfig[userModel, users.User]{
		RecordSchema: map[string]schema.Attribute{
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
		ResultAttributeName: "users",
		PrettyAttributeName: "users",
		FlattenAPIModel: func(ctx context.Context, apiObject users.User) (state *userModel, diags diag.Diagnostics) {
			state = new(userModel)

			state.ID = types.StringValue(apiObject.ID)
			state.OrganizationID = types.StringValue(apiObject.OrganizationID)
			state.FullName = types.StringValue(apiObject.FullName)
			state.Email = types.StringValue(apiObject.Email)
			state.IsAdmin = types.BoolValue(apiObject.IsAdmin)
			state.TimeCreated = types.StringValue(apiObject.TimeCreated.UTC().Format(time.RFC3339))

			if apiObject.LastLogin != nil {
				state.LastLogin = types.StringValue(apiObject.LastLogin.UTC().Format(time.RFC3339))
			} else {
				state.LastLogin = types.StringNull()
			}

			return
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]users.User, error) {
			users, _, err := client.Users.ListUsers(ctx)
			return users, err
		},
		Description: "Get a list of all users in your BastionZero organization.",
	})
}
