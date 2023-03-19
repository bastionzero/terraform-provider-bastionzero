package bzdatasource

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TFModelWithID interface {
	GetID() types.String
}

// SingleDataSourceConfig is the configuration for a single item data source. It
// represents the schema and operations needed to create the single data source.
type SingleDataSourceConfig[TFModel TFModelWithID, APIModel any] struct {
	// RecordSchema is the TF schema that models a single instance of the API
	// object. Required. Schema must contain a required attribute with name
	// "id".
	RecordSchema map[string]schema.Attribute

	// The name of the attribute in the data source definition name. Cannot be
	// the empty string.
	ResultAttributeName string

	// PrettyAttributeName is the name of the attribute used for logging and
	// documentation purposes. Cannot be the empty string.
	PrettyAttributeName string

	// Given a model returned from the GetAPIModel function, flatten the API
	// model to a TF model.
	FlattenAPIModel func(ctx context.Context, apiObject *APIModel) (*TFModel, diag.Diagnostics)

	// GetAPIModel returns a single API object, on which the data source should
	// expose, given the ID.
	GetAPIModel func(ctx context.Context, client *bastionzero.Client, id string) (*APIModel, error)

	// Description is passed as the data source schema's Description field
	// during construction.
	Description string

	// MarkdownDescription is passed as the data source schema's
	// MarkdownDescription field during construction.
	MarkdownDescription string

	// DeprecationMessage is passed as the data source schema's
	// DepcrecationMessage field during construction.
	DeprecationMessage string
}

// Returns a new single data source given the specified configuration. The
// function panics if the config is invalid. A single data source abstracts
// calling a GET BastionZero API endpoint that takes an ID and returns a single
// object.
func NewSingleDataSource[TFModel TFModelWithID, APIModel any](config *SingleDataSourceConfig[TFModel, APIModel]) datasource.DataSourceWithConfigure {
	if config.RecordSchema == nil {
		panic("RecordSchema cannot be nil")
	}
	if val, ok := config.RecordSchema["id"]; !ok || !val.IsRequired() {
		panic("RecordSchema must contain attribute with name \"id\" and it must have \"Required\" set to true")
	}
	if config.ResultAttributeName == "" {
		panic("ResultAttributeName cannot be empty")
	}
	if config.PrettyAttributeName == "" {
		panic("PrettyAttributeName cannot be empty")
	}

	t := struct{ protoDataSource }{}
	t.metadataFunc = func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
		resp.TypeName = req.ProviderTypeName + fmt.Sprintf("_%s", config.ResultAttributeName)
	}
	t.schemaFunc = func(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
		resp.Schema = schema.Schema{
			Description:         config.Description,
			MarkdownDescription: config.MarkdownDescription,
			DeprecationMessage:  config.DeprecationMessage,
			Attributes:          config.RecordSchema,
		}
	}
	t.configureFunc = func(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

		t.client = client
	}
	t.readFunc = func(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
		var T TFModel
		// Read Terraform configuration data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, &T)...)
		if resp.Diagnostics.HasError() {
			return
		}
		id := T.GetID()
		ctx = tflog.SetField(ctx, fmt.Sprintf("%s_id", config.ResultAttributeName), id.ValueString())

		// Query BastionZero for API object
		tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.PrettyAttributeName))
		apiObject, err := config.GetAPIModel(ctx, t.client, id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName))

		// Convert to TFModel
		tfModel, diags := config.FlattenAPIModel(ctx, apiObject)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
	}

	return &t
}
