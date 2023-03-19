package bastionzero

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/organization"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/targetconnect"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/serviceaccount"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target/bzerotarget"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/user"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure BastionZeroProvider satisfies various provider interfaces.
var _ provider.Provider = &BastionZeroProvider{}

// BastionZeroProvider defines the provider implementation.
type BastionZeroProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// bastionzeroProviderModel describes the provider data model.
type bastionzeroProviderModel struct {
	Host      types.String `tfsdk:"api_endpoint"`
	APISecret types.String `tfsdk:"api_secret"`
}

func (p *BastionZeroProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bastionzero"
	resp.Version = p.version
}

func (p *BastionZeroProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The BastionZero provider is used to interact with select APIs provided by BastionZero. " +
			"The provider needs to be configured with an API secret before it can be used (please see the example below).\n\n" +
			"Use the navigation to the left to read about the available resources and data sources.",
		Description: "Provider for the BastionZero API",
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				Description: fmt.Sprintf("This can be used to override the base URL for BastionZero API requests (Defaults to the value of the BASTIONZERO_HOST environment variable or %s if unset)."+
					"Typical users of this provider should not set this value.", bastionzero.DefaultBaseURL),
				Optional: true,
			},
			"api_secret": schema.StringAttribute{
				Description: "API secret used to authenticate API requests sent to BastionZero. This can also be specified using the BASTIONZERO_API_SECRET environment variable.",
				Sensitive:   true,
				Optional:    true,
			},
		},
	}
}

func (p *BastionZeroProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	ctx = tflog.SetField(ctx, "bastionzero_provider_version", p.version)
	tflog.Info(ctx, "Configuring BastionZero client")

	var config bastionzeroProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the attributes,
	// it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_endpoint"),
			"Unknown BastionZero API Host",
			"The provider cannot create the BastionZero API client as there is an unknown configuration value for the BastionZero API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BASTIONZERO_HOST environment variable.",
		)
	}

	if config.APISecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Unknown BastionZero API Secret",
			"The provider cannot create the BastionZero API client as there is an unknown configuration value for the BastionZero API secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BASTIONZERO_API_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform
	// configuration value if set.

	host := os.Getenv("BASTIONZERO_HOST")
	apiSecret := os.Getenv("BASTIONZERO_API_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.APISecret.IsNull() {
		apiSecret = config.APISecret.ValueString()
	}

	// If any of the expected configurations are missing, return errors with
	// provider-specific guidance.

	if apiSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Missing BastionZero API Secret",
			"The provider cannot create the BastionZero API client as there is a missing or empty value for the BastionZero API secret. "+
				"Set the api_secret value in the configuration or use the BASTIONZERO_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Always include user agent header
	opts := []bastionzero.ClientOpt{bastionzero.WithUserAgent("terraform-provider-bastionzero/" + p.version)}
	// If custom host specified, configure client with base URL
	if host != "" {
		opts = append(opts, bastionzero.WithBaseURL(host))
	}

	ctx = tflog.SetField(ctx, "bastionzero_host", host)
	ctx = tflog.SetField(ctx, "bastionzero_api_secret", apiSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "bastionzero_api_secret")

	tflog.Debug(ctx, "Creating BastionZero client")

	// f, err := os.OpenFile("foo.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// reqClient := reqc.C()
	// reqClient.EnableDumpAllTo(f)
	// reqClient.EnableDebugLog()

	// c := &http.Client{Transport: reqClient.GetTransport()}

	// Create a new BastionZero client using the configuration values
	client, err := bastionzero.NewFromAPISecret(http.DefaultClient, apiSecret, opts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create BastionZero API Client",
			"An unexpected error occurred when creating the BastionZero API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"BastionZero Client Error: "+err.Error(),
		)
		return
	}

	// Make the BastionZero client available during DataSource and Resource type
	// Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured BastionZero client", map[string]any{"success": true})
}

func (p *BastionZeroProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		environment.NewEnvironmentResource,
		targetconnect.NewTargetConnectPolicyResource,
	}
}

func (p *BastionZeroProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		user.NewUserDataSource,
		user.NewUsersDataSource,
		organization.NewGroupsDataSource,
		serviceaccount.NewServiceAccountsDataSource,
		environment.NewEnvironmentDataSource,
		environment.NewEnvironmentsDataSource,
		bzerotarget.NewBzeroTargetsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BastionZeroProvider{
			version: version,
		}
	}
}
