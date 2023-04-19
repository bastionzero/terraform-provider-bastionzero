package bzdatasource

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Anonymous interface implementation via prototype struct and embedding
//
// Source: https://stackoverflow.com/a/31362378

// Ensure prototype implements data source framework interface
var _ datasource.DataSourceWithConfigure = &protoDataSource{}

type protoDataSource struct {
	client *bastionzero.Client

	metadataFunc  func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	schemaFunc    func(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	configureFunc func(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
	readFunc      func(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

func (p *protoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	p.metadataFunc(ctx, req, resp)
}
func (p *protoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	p.schemaFunc(ctx, req, resp)
}
func (p *protoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	p.configureFunc(ctx, req, resp)
}
func (p *protoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	p.readFunc(ctx, req, resp)
}
