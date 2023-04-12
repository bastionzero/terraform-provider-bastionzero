package bastionzero

import (
	"context"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func (p *BastionZeroProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bastionzero"
	resp.Version = p.version
}

func (p *BastionZeroProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for the BastionZero API",
	}
}

func (p *BastionZeroProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *BastionZeroProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		environment.NewEnvironmentResource,
	}
}

func (p *BastionZeroProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// New creates a BastionZero Terraform provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BastionZeroProvider{
			version: version,
		}
	}
}
