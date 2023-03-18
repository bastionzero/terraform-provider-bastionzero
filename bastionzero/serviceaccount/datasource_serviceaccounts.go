package serviceaccount

import (
	"context"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// serviceAccountModel maps service account schema data.
type serviceAccountModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Email          types.String `tfsdk:"email"`
	ExternalID     types.String `tfsdk:"external_id"`
	JwksURL        types.String `tfsdk:"jwks_url"`
	JwksURLPattern types.String `tfsdk:"jwks_url_pattern"`
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
	TimeCreated    types.String `tfsdk:"time_created"`
	LastLogin      types.String `tfsdk:"last_login"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func NewServiceAccountsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[serviceAccountModel, serviceaccounts.ServiceAccount]{
		RecordSchema: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The service account's unique ID.",
			},
			"organization_id": schema.StringAttribute{
				Computed:    true,
				Description: "The service account's organization's ID.",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "The service account's email address.",
			},
			"external_id": schema.StringAttribute{
				Computed:    true,
				Description: "The service account's unique per service provider identifier provided by the user during creation.",
			},
			"jwks_url": schema.StringAttribute{
				Computed:    true,
				Description: "The service account's publicly available JWKS URL that provides the public key that can be used to verify the tokens signed by the private key of this service account.",
			},
			"jwks_url_pattern": schema.StringAttribute{
				Computed:    true,
				Description: " A URL pattern that all service accounts of the same service account provider follow in their JWKS URL.",
			},
			"is_admin": schema.BoolAttribute{
				Computed:    true,
				Description: "If true, the service account is an administrator. False otherwise.",
			},
			"time_created": schema.StringAttribute{
				Computed:    true,
				Description: "The time this service account was created in BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
			},
			"last_login": schema.StringAttribute{
				Computed:    true,
				Description: "The time this service account last logged into BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
				Optional:    true,
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the subject that created this service account.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "If true, the service account is currently enabled. False otherwise.",
			},
		},
		ResultAttributeName: "service_accounts",
		PrettyAttributeName: "service accounts",
		FlattenAPIModel: func(ctx context.Context, apiObject serviceaccounts.ServiceAccount) (state *serviceAccountModel, diags diag.Diagnostics) {
			state = new(serviceAccountModel)
			state.ID = types.StringValue(apiObject.ID)
			state.OrganizationID = types.StringValue(apiObject.OrganizationID)
			state.Email = types.StringValue(apiObject.Email)
			state.ExternalID = types.StringValue(apiObject.ExternalID)
			state.JwksURL = types.StringValue(apiObject.JwksURL)
			state.JwksURLPattern = types.StringValue(apiObject.JwksURLPattern)
			state.IsAdmin = types.BoolValue(apiObject.IsAdmin)
			state.TimeCreated = types.StringValue(apiObject.TimeCreated.UTC().Format(time.RFC3339))
			state.CreatedBy = types.StringValue(apiObject.CreatedBy)
			state.Enabled = types.BoolValue(apiObject.Enabled)

			if apiObject.LastLogin != nil {
				state.LastLogin = types.StringValue(apiObject.LastLogin.UTC().Format(time.RFC3339))
			} else {
				state.LastLogin = types.StringNull()
			}

			return
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]serviceaccounts.ServiceAccount, error) {
			serviceAccounts, _, err := client.ServiceAccounts.ListServiceAccounts(ctx)
			return serviceAccounts, err
		},
		Description: "Get a list of all service accounts in your BastionZero organization.",
	})
}
