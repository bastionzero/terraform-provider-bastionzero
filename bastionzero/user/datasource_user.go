package user

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewUserDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[userModel, users.User]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[userModel, users.User]{
				RecordSchema:        makeUserDataSourceSchema(true),
				MetadataTypeName:    "user",
				PrettyAttributeName: "user",
				FlattenAPIModel: func(ctx context.Context, apiObject *users.User, state *userModel) (diags diag.Diagnostics) {
					setUserAttributes(ctx, state, apiObject)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel userModel, client *bastionzero.Client) (*users.User, error) {
					user, _, err := client.Users.GetUser(ctx, tfModel.ID.ValueString())
					return user, err
				},
				MarkdownDescription: "Get information on a user in your BastionZero organization.",
			},
		},
	)
}

// Provide the user's unique ID or email address in the \"id\" field.
