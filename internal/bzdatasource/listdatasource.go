package bzdatasource

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jinzhu/copier"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"golang.org/x/exp/maps"
)

type BaseListDataSourceConfig[TFComputedModel any, APIModel any] struct {
	// RecordSchema is the TF schema that models a single instance of the API
	// object. Required.
	RecordSchema map[string]schema.Attribute

	// The name of the attribute in the data source through which to expose a
	// list of results. Cannot be the empty string.
	ResultAttributeName string

	// The suffix name to use for the data source name. Optional. If not set,
	// then PrettyAttributeName is used.
	MetadataTypeName string

	// PrettyAttributeName is the name of the attribute used for logging and
	// documentation purposes. Cannot be the empty string.
	PrettyAttributeName string

	// Given a model returned from the ListAPIModels function, flatten the API
	// model to a TF model.
	FlattenAPIModel func(ctx context.Context, apiObject *APIModel) (*TFComputedModel, diag.Diagnostics)

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

// ListDataSourceConfig is the configuration for a list data source. It
// represents the schema and operations needed to create the list data source.
type ListDataSourceConfig[TFComputedModel any, APIModel any] struct {
	*BaseListDataSourceConfig[TFComputedModel, APIModel]

	// ListAPIModels returns all of the API models on which the data source
	// should expose.
	ListAPIModels func(ctx context.Context, client *bastionzero.Client) ([]APIModel, error)
}

// Returns a new list data source given the specified configuration. The
// function panics if the config is invalid. A list data source abstracts
// calling a GET BastionZero API endpoint that returns a list of objects.
func NewListDataSource[TFModel any, APIModel any](config *ListDataSourceConfig[TFModel, APIModel]) datasource.DataSourceWithConfigure {
	if config.RecordSchema == nil {
		panic("RecordSchema cannot be nil")
	}
	if config.ResultAttributeName == "" {
		panic("ResultAttributeName cannot be empty")
	}
	if config.PrettyAttributeName == "" {
		panic("PrettyAttributeName cannot be empty")
	}

	if config.MetadataTypeName == "" {
		config.MetadataTypeName = config.ResultAttributeName
	}

	t := struct{ protoDataSource }{}
	t.metadataFunc = func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
		resp.TypeName = req.ProviderTypeName + fmt.Sprintf("_%s", config.MetadataTypeName)
	}
	t.schemaFunc = func(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
		resp.Schema = schema.Schema{
			Description:         config.Description,
			MarkdownDescription: config.MarkdownDescription,
			DeprecationMessage:  config.DeprecationMessage,
			Attributes: map[string]schema.Attribute{
				config.ResultAttributeName: schema.ListNestedAttribute{
					Description: fmt.Sprintf("List of %s.", config.PrettyAttributeName),
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
		tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.PrettyAttributeName))
		apiObjects, err := config.ListAPIModels(ctx, t.client)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to list %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName), map[string]any{fmt.Sprintf("num_%s", config.ResultAttributeName): len(apiObjects)})

		// Map response body to model
		for _, apiObj := range apiObjects {
			tfModel, diags := config.FlattenAPIModel(ctx, &apiObj)
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

// ListDataSourceConfig is the configuration for a list data source. It
// represents the schema and operations needed to create the list data source.
type ListDataSourceConfigWithUserParameters[TFComputedModel any, TFUserParamsModel any, APIModel any] struct {
	*BaseListDataSourceConfig[TFComputedModel, APIModel]

	// UserParamsRecordSchema is the TF schema that models additional user
	// parameters that are passed to ListAPIModels. Required.
	UserParamsRecordSchema map[string]schema.Attribute

	// ListAPIModels returns all of the API models on which the data source
	// should expose.
	ListAPIModels func(ctx context.Context, tfUserParamsModel TFUserParamsModel, client *bastionzero.Client) ([]APIModel, error)
}

// Returns a new list data source given the specified configuration. The
// function panics if the config is invalid. A list data source abstracts
// calling a GET BastionZero API endpoint that returns a list of objects.
func NewListDataSourceWithUserParameters[TFModel any, TFUserParamsModel any, APIModel any](config *ListDataSourceConfigWithUserParameters[TFModel, TFUserParamsModel, APIModel]) datasource.DataSourceWithConfigure {
	if config.RecordSchema == nil {
		panic("RecordSchema cannot be nil")
	}
	if config.UserParamsRecordSchema == nil {
		panic("UserParamsRecordSchema cannot be nil")
	}
	if config.ResultAttributeName == "" {
		panic("ResultAttributeName cannot be empty")
	}
	if config.PrettyAttributeName == "" {
		panic("PrettyAttributeName cannot be empty")
	}
	if _, ok := config.UserParamsRecordSchema[config.ResultAttributeName]; ok {
		panic(fmt.Sprintf("UserParamsRecordSchema cannot have attribute with name %v", config.ResultAttributeName))
	}

	if config.MetadataTypeName == "" {
		config.MetadataTypeName = config.ResultAttributeName
	}

	t := struct{ protoDataSource }{}
	t.metadataFunc = func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
		resp.TypeName = req.ProviderTypeName + fmt.Sprintf("_%s", config.MetadataTypeName)
	}
	t.schemaFunc = func(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
		attributes := map[string]schema.Attribute{
			config.ResultAttributeName: schema.ListNestedAttribute{
				Description: fmt.Sprintf("List of %s.", config.PrettyAttributeName),
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: config.RecordSchema,
				},
			},
		}
		maps.Copy(attributes, config.UserParamsRecordSchema)

		resp.Schema = schema.Schema{
			Description:         config.Description,
			MarkdownDescription: config.MarkdownDescription,
			DeprecationMessage:  config.DeprecationMessage,
			Attributes:          attributes,
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
		var userParamsModel TFUserParamsModel
		stateScaffold := struct{ Records []TFModel }{}

		mergedModelBuilder := dynamicstruct.MergeStructs(stateScaffold, userParamsModel)

		// Add required tag so we can get and write to TF state
		recordsField := mergedModelBuilder.GetField("Records")
		recordsField.SetTag(fmt.Sprintf(`tfsdk:"%s"`, config.ResultAttributeName))

		mergedModel := mergedModelBuilder.Build().New()

		// Read Terraform configuration data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, mergedModel)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Using reflection, populate userParamsModel
		err := copier.Copy(&userParamsModel, mergedModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error during conversion to model with user parameters",
				fmt.Sprintf("Got: %v during conversion. Please report this issue to the provider developers.", err.Error()),
			)
			return
		}

		// Query BastionZero for list of API objects
		tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.PrettyAttributeName))
		apiObjects, err := config.ListAPIModels(ctx, userParamsModel, t.client)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to list %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName), map[string]any{fmt.Sprintf("num_%s", config.ResultAttributeName): len(apiObjects)})

		// Map response body to model
		for _, apiObj := range apiObjects {
			tfModel, diags := config.FlattenAPIModel(ctx, &apiObj)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				// Return early because something went wrong
				return
			}

			stateScaffold.Records = append(stateScaffold.Records, *tfModel)
		}

		// Using reflection, copy values from records back into the merged model
		// that is expected to be stored in TF state
		err = copier.Copy(mergedModel, stateScaffold)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error during conversion to merged model",
				fmt.Sprintf("Got: %v during conversion. Please report this issue to the provider developers.", err.Error()),
			)
			return
		}

		// // Dynamically set the TF state struct's tfsdk tag
		// //
		// // Source: https://stackoverflow.com/a/62486560
		// value := reflect.ValueOf(stateScaffold)
		// stateScaffoldType := value.Type()
		// sf := make([]reflect.StructField, 0)
		// sf = append(sf, stateScaffoldType.Field(0))
		// sf[0].Tag = reflect.StructTag(fmt.Sprintf(`tfsdk:"%s"`, config.ResultAttributeName))
		// newType := reflect.StructOf(sf)
		// newValue := value.Convert(newType)

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, mergedModel)...)
	}

	return &t
}
