package listdatasource

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ResourceConfig is the configuration for a list data source. It represents the
// schema and operations needed to create the list data source.
type ResourceConfig[TFModel any, APIModel any] struct {
	// The schema for a single instance of the resource.
	RecordSchema map[string]schema.Attribute

	// The name of the attribute in the resource through which to expose
	// results.
	ResultAttributeName string

	// Given a model returned from the ListAPIModels function, flatten the API
	// model to a TF model.
	FlattenAPIModel func(ctx context.Context, apiObject APIModel) (*TFModel, diag.Diagnostics)

	// ListAPIModels returns all of the API models on which the data source
	// should expose.
	ListAPIModels func(ctx context.Context, client *bastionzero.Client) ([]APIModel, error)

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

// Returns a new list data source given the specified configuration.
func NewListDataSource[TFModel any, APIModel any](config *ResourceConfig[TFModel, APIModel]) datasource.DataSourceWithConfigure {
	t := struct{ protoDataSource }{}
	t.metadataFunc = func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
		resp.TypeName = req.ProviderTypeName + fmt.Sprintf("_%s", config.ResultAttributeName)
	}
	t.schemaFunc = func(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
		resp.Schema = schema.Schema{
			Description:         config.Description,
			MarkdownDescription: config.MarkdownDescription,
			DeprecationMessage:  config.DeprecationMessage,
			Attributes: map[string]schema.Attribute{
				config.ResultAttributeName: schema.ListNestedAttribute{
					Description: fmt.Sprintf("List of %s.", config.ResultAttributeName),
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: config.RecordSchema,
					},
				},
			},
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
		stateScaffold := struct{ Records []TFModel }{}

		// Query BastionZero for list of API objects
		tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.ResultAttributeName))
		apiObjects, err := config.ListAPIModels(ctx, t.client)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to list %s", config.ResultAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.ResultAttributeName), map[string]any{fmt.Sprintf("num_%s", config.ResultAttributeName): len(apiObjects)})

		// Map response body to model
		for _, apiObj := range apiObjects {
			tfModel, diags := config.FlattenAPIModel(ctx, apiObj)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				// Return early because something went wrong
				return
			}

			stateScaffold.Records = append(stateScaffold.Records, *tfModel)
		}

		// Dynamically set the TF state struct's tfsdk tag
		//
		// Source: https://stackoverflow.com/a/62486560
		value := reflect.ValueOf(stateScaffold)
		stateScaffoldType := value.Type()
		sf := make([]reflect.StructField, 0)
		sf = append(sf, stateScaffoldType.Field(0))
		sf[0].Tag = reflect.StructTag(fmt.Sprintf(`tfsdk:"%s"`, config.ResultAttributeName))
		newType := reflect.StructOf(sf)
		newValue := value.Convert(newType)

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, newValue.Interface())...)
	}

	return &t
}

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

func (t *protoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	t.metadataFunc(ctx, req, resp)
}
func (t *protoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.schemaFunc(ctx, req, resp)
}
func (t *protoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	t.configureFunc(ctx, req, resp)
}
func (t *protoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	t.readFunc(ctx, req, resp)
}
