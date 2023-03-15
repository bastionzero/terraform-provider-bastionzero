package serviceaccount

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
var _ datasource.DataSource = &serviceAccountsDataSource{}

func NewServiceAccountsDataSource() datasource.DataSource {
	return &serviceAccountsDataSource{}
}

// serviceAccountsDataSource is the data source implementation.
type serviceAccountsDataSource struct {
	client *bastionzero.Client
}

// serviceAccountsDataSourceModel describes the service accounts data source
// data model.
type serviceAccountsDataSourceModel struct {
	ServiceAccounts []serviceAccountModel `tfsdk:"service_accounts"`
}

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

// Metadata returns the service accounts data source type name.
func (d *serviceAccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_accounts"
}

// Schema defines the schema for the service accounts data source.
func (d *serviceAccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a list of all service accounts in your BastionZero organization.",
		Attributes: map[string]schema.Attribute{
			"service_accounts": schema.ListNestedAttribute{
				Description: "List of service accounts in your organization.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
				},
			},
		},
	}
}

// Configure adds the provider configured BastionZero API client to the data
// source.
func (d *serviceAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serviceAccountsDataSourceModel

	// Query BastionZero for service accounts
	tflog.Debug(ctx, "Querying for service accounts")
	serviceAccounts, _, err := d.client.ServiceAccounts.ListServiceAccounts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list service accounts",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Queried for service accounts", map[string]any{"num_service_accounts": len(serviceAccounts)})

	// Map response body to model
	for _, serviceAccount := range serviceAccounts {
		serviceAccountState := serviceAccountModel{
			ID:             types.StringValue(serviceAccount.ID),
			OrganizationID: types.StringValue(serviceAccount.OrganizationID),
			Email:          types.StringValue(serviceAccount.Email),
			ExternalID:     types.StringValue(serviceAccount.ExternalID),
			JwksURL:        types.StringValue(serviceAccount.JwksURL),
			JwksURLPattern: types.StringValue(serviceAccount.JwksURLPattern),
			IsAdmin:        types.BoolValue(serviceAccount.IsAdmin),
			TimeCreated:    types.StringValue(serviceAccount.TimeCreated.UTC().Format(time.RFC3339)),
			CreatedBy:      types.StringValue(serviceAccount.CreatedBy),
			Enabled:        types.BoolValue(serviceAccount.Enabled),
		}

		if serviceAccount.LastLogin != nil {
			serviceAccountState.LastLogin = types.StringValue(serviceAccount.LastLogin.UTC().Format(time.RFC3339))
		} else {
			serviceAccountState.LastLogin = types.StringNull()
		}

		state.ServiceAccounts = append(state.ServiceAccounts, serviceAccountState)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
