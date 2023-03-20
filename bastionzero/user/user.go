package user

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// userModel maps user schema data.
type userModel struct {
	ID             types.String `tfsdk:"id"`
	Type           types.String `tfsdk:"type"`
	OrganizationID types.String `tfsdk:"organization_id"`
	FullName       types.String `tfsdk:"full_name"`
	Email          types.String `tfsdk:"email"`
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
	TimeCreated    types.String `tfsdk:"time_created"`
	LastLogin      types.String `tfsdk:"last_login"`
}

// setUserAttributes populates the TF schema data from a user API object.
func setUserAttributes(ctx context.Context, schema *userModel, user *users.User) {
	schema.ID = types.StringValue(user.ID)
	schema.Type = types.StringValue(string(subjecttype.User))
	schema.OrganizationID = types.StringValue(user.OrganizationID)
	schema.FullName = types.StringValue(user.FullName)
	schema.Email = types.StringValue(user.Email)
	schema.IsAdmin = types.BoolValue(user.IsAdmin)
	schema.TimeCreated = types.StringValue(user.TimeCreated.UTC().Format(time.RFC3339))

	if user.LastLogin != nil {
		schema.LastLogin = types.StringValue(user.LastLogin.UTC().Format(time.RFC3339))
	} else {
		schema.LastLogin = types.StringNull()
	}
}

func makeUserDataSourceSchema(withRequiredID bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    !withRequiredID,
			Required:    withRequiredID,
			Description: "The user's unique ID.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The subject's type (constant value \"%s\").", subjecttype.User),
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
			Description: "The time this user last logged into BastionZero formatted as a UTC timestamp string in RFC 3339 format. Null if the user has never logged in.",
		},
	}
}
