package organization

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/organization"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// groupModel maps group schema data.
type groupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewGroupsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[groupModel, organization.Group]{
		RecordSchema: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The group's unique ID, as specified by the Identity Provider in which it is configured.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The group's name.",
			},
		},
		ResultAttributeName: "groups",
		PrettyAttributeName: "groups",
		FlattenAPIModel: func(ctx context.Context, apiObject *organization.Group) (state *groupModel, diags diag.Diagnostics) {
			state = new(groupModel)
			state.ID = types.StringValue(apiObject.ID)
			state.Name = types.StringValue(apiObject.Name)

			return
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]organization.Group, error) {
			groups, _, err := client.Organization.ListGroups(ctx)
			return groups, err
		},
		Description: "Get a list of all groups in your BastionZero organization.",
	})
}
