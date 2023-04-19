package bzdatasource

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jinzhu/copier"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"golang.org/x/exp/maps"

	"github.com/google/uuid"
)

// TFComputedModel is a struct that models a collection of TF schema attributes
// that are Computed (Optional and Required should both be set to false).
type TFComputedModel = interface{}

// BaseListDataSourceConfig contains common options used for creating any type
// of ListDataSource data source.
type BaseListDataSourceConfig[T TFComputedModel, T2 APIModel] struct {
	// RecordSchema is the TF schema that models a single instance of the API
	// object. There should be a key for each field defined in TFComputedModel.
	// Required.
	RecordSchema map[string]schema.Attribute

	// ResultAttributeName is the name of the TF attribute in the data source
	// through which to expose a list of results. Cannot be the empty string.
	ResultAttributeName string

	// MetadataTypeName is the suffix to use for the name of the data source.
	// Optional. If not set, then ResultAttributeName is used.
	MetadataTypeName string

	// PrettyAttributeName is a descriptive name used for logging and
	// documentation purposes. Cannot be the empty string.
	PrettyAttributeName string

	// FlattenAPIModel takes a model returned from the ListAPIModels function
	// and converts it to a TF model
	FlattenAPIModel func(ctx context.Context, apiObject *T2) (*T, diag.Diagnostics)

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

// Validate checks for errors in the base list data source config and fills in
// required values that were not set
func (c *BaseListDataSourceConfig[T, T2]) Validate() error {
	if c.RecordSchema == nil {
		return fmt.Errorf("RecordSchema cannot be nil")
	}
	if c.ResultAttributeName == "" {
		return fmt.Errorf("ResultAttributeName cannot be empty")
	}
	if c.PrettyAttributeName == "" {
		return fmt.Errorf("PrettyAttributeName cannot be empty")
	}
	if c.MetadataTypeName == "" {
		c.MetadataTypeName = c.ResultAttributeName
	}

	return nil
}

// ListDataSource is a data source that calls a BastionZero API endpoint which
// returns a list of objects. It abstracts common, boilerplate code that
// typically accompanies a data source that exposes a list of items.
//
// If the API endpoint takes in additional parameters, and you wish to expose
// these to the practitioner for configuration in the TF schema, then use
// bzdatasource.ListDataSourceWithPractitionerParameters instead.
type ListDataSource datasource.DataSourceWithConfigure

// ListDataSourceConfig is the configuration for a list data source. It
// represents the schema and operations needed to create the data source.
type ListDataSourceConfig[T TFComputedModel, T2 APIModel] struct {
	*BaseListDataSourceConfig[T, T2]

	// ListAPIModels returns all of the API models on which the data source
	// should expose.
	ListAPIModels func(ctx context.Context, client *bastionzero.Client) ([]T2, error)
}

// NewListDataSource creates a ListDataSource. The function panics if the config
// is invalid.
func NewListDataSource[T TFComputedModel, T2 APIModel](config *ListDataSourceConfig[T, T2]) ListDataSource {
	if err := config.Validate(); err != nil {
		panic(err)
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
				// A list data source exposes a single attribute; a list of
				// TFComputedModel objects.
				config.ResultAttributeName: schema.ListNestedAttribute{
					Description: fmt.Sprintf("List of %s.", config.PrettyAttributeName),
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: config.RecordSchema,
					},
				},
				// Dummy "id" attribute. Required in order to test this data source.
				//
				// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84
				// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
				// Source: https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
				"id": schema.StringAttribute{
					Computed:           true,
					Description:        "Deprecated. Do not depend on this attribute. This attribute will be removed in the future.",
					DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
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
		stateScaffold := struct {
			Records []T
			Id      *string `tfsdk:"id"`
		}{}

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
		models, diags := mapAPIModelsToTFModels(ctx, config.BaseListDataSourceConfig, apiObjects)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		stateScaffold.Records = models
		stateScaffold.Id = bastionzero.PtrTo(uuid.New().String())

		// Dynamically set the TF state struct's tfsdk tag
		modelBuilder := dynamicstruct.ExtendStruct(stateScaffold)
		recordsField := modelBuilder.GetField("Records")
		recordsField.SetTag(fmt.Sprintf(`tfsdk:"%s"`, config.ResultAttributeName))
		model := modelBuilder.Build().New()

		err = copier.Copy(model, stateScaffold)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error during conversion to model",
				fmt.Sprintf("Got: %v during conversion. Please report this issue to the provider developers.", err.Error()),
			)
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	}

	return &t
}

// ListDataSourceWithPractitionerParameters is a data source that calls a
// BastionZero API endpoint which returns a list of objects. It abstracts
// common, boilerplate code that typically accompanies a data source that
// exposes a list of items. Additionally, it provides support for practitioner
// provided attributes (either Required = true or Optional = true) that can be
// accessed before calling the BastionZero API endpoint.
//
// If there are no extra practitioner parameters required to call the API
// endpoint, then use bzdatasource.ListDataSource instead.
type ListDataSourceWithPractitionerParameters datasource.DataSourceWithConfigure

// TFNonComputedModel is a struct that models a collection of TF schema
// attributes that are not Computed (Computed should be set to false. Required
// or Optional can be set to True).
type TFNonComputedModel = interface{}

// ListDataSourceWithPractitionerParametersConfig is the configuration for a
// list data source with practitioner parameters. It represents the schema and
// operations needed to create the data source.
type ListDataSourceWithPractitionerParametersConfig[T TFComputedModel, T2 TFNonComputedModel, T3 APIModel] struct {
	*BaseListDataSourceConfig[T, T3]

	// PractitionerParamsRecordSchema is the TF schema that models additional
	// user parameters that are passed to ListAPIModels. Required.
	PractitionerParamsRecordSchema map[string]schema.Attribute

	// ListAPIModels returns all of the API models on which the data source
	// should expose. practitionerParams are the practitioner parameters
	// retrieved from the TF schema.
	ListAPIModels func(ctx context.Context, practitionerParams T2, client *bastionzero.Client) ([]T3, error)
}

// NewListDataSourceWithPractitionerParameters creates a
// ListDataSourceWithPractitionerParameters. The function panics if the config
// is invalid.
func NewListDataSourceWithPractitionerParameters[T TFComputedModel, T2 TFNonComputedModel, T3 APIModel](config *ListDataSourceWithPractitionerParametersConfig[T, T2, T3]) ListDataSourceWithPractitionerParameters {
	if err := config.Validate(); err != nil {
		panic(err)
	}
	if config.PractitionerParamsRecordSchema == nil {
		panic("PractitionerParamsRecordSchema cannot be nil")
	}
	if _, ok := config.PractitionerParamsRecordSchema[config.ResultAttributeName]; ok {
		panic(fmt.Sprintf("PractitionerParamsRecordSchema cannot have attribute with name %v", config.ResultAttributeName))
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
			// Dummy "id" attribute. Required in order to test this data source.
			//
			// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84
			// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
			// Source: https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed:           true,
				Description:        "Deprecated. Do not depend on this attribute. This attribute will be removed in the future.",
				DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
			},
		}
		// Add extra practitioner parameters
		maps.Copy(attributes, config.PractitionerParamsRecordSchema)

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
		var userParamsModel T2
		stateScaffold := struct {
			Records []T
			Id      *string `tfsdk:"id"`
		}{}

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
		models, diags := mapAPIModelsToTFModels(ctx, config.BaseListDataSourceConfig, apiObjects)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		stateScaffold.Records = models
		stateScaffold.Id = bastionzero.PtrTo(uuid.New().String())

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

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, mergedModel)...)
	}

	return &t
}

// mapAPIModelsToTFModels converts a list of BastionZero API models to a list of
// TF models.
func mapAPIModelsToTFModels[T TFComputedModel, T2 APIModel](ctx context.Context, config *BaseListDataSourceConfig[T, T2], apiObjects []T2) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics
	models := make([]T, 0)
	for _, apiObj := range apiObjects {
		tfModel, diagsFlatten := config.FlattenAPIModel(ctx, &apiObj)
		diags.Append(diagsFlatten...)
		if diagsFlatten.HasError() {
			// Return early because something went wrong
			return models, diags
		}

		models = append(models, *tfModel)
	}

	return models, diags
}
