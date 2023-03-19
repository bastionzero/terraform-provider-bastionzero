package user

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/users"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewUsersDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[userModel, users.User]{
		RecordSchema:        makeUserDataSourceSchema(false),
		ResultAttributeName: "users",
		PrettyAttributeName: "users",
		FlattenAPIModel: func(ctx context.Context, apiObject *users.User) (state *userModel, diags diag.Diagnostics) {
			state = new(userModel)
			setUserAttributes(ctx, state, apiObject)
			return
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]users.User, error) {
			users, _, err := client.Users.ListUsers(ctx)
			return users, err
		},
		Description: "Get a list of all users in your BastionZero organization.",
	})
}
