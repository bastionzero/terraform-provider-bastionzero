package bzdatasource

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jinzhu/copier"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

// TODO: Refactor FlattenAPIModel and GetAPIModel functions to take in req +
// resp similar to the Terraform plugin framework, so that we can add new fields
// without having to refactor every instance of these structs

// TODO: FlattenAPIModel: Potentially consider taking pointer to state that was
// read previously instead of asking for new value to be returned and passing
// copy of current value read

type BaseSingleDataSourceConfig[TFModel any, APIModel any] struct {
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
	FlattenAPIModel func(ctx context.Context, apiObject *APIModel, tfModel TFModel) (*TFModel, diag.Diagnostics)

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

// SingleDataSourceConfig is the configuration for a single item data source. It
// represents the schema and operations needed to create the single data source.
type SingleDataSourceConfig[TFModel any, APIModel any] struct {
	*BaseSingleDataSourceConfig[TFModel, APIModel]

	// GetAPIModel returns a single API object on which the data source should
	// expose.
	GetAPIModel func(ctx context.Context, tfModel TFModel, client *bastionzero.Client) (*APIModel, error)
}

// Returns a new single data source given the specified configuration. The
// function panics if the config is invalid. A single data source abstracts
// calling a BastionZero API endpoint that returns a single object.
func NewSingleDataSource[TFModel any, APIModel any](config *SingleDataSourceConfig[TFModel, APIModel]) datasource.DataSourceWithConfigure {
	if config.RecordSchema == nil {
		panic("RecordSchema cannot be nil")
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

		// Query BastionZero for API object
		tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.PrettyAttributeName))
		apiObject, err := config.GetAPIModel(ctx, T, t.client)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName))

		// Convert to TFModel
		tfModel, diags := config.FlattenAPIModel(ctx, apiObject, T)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, tfModel)...)
	}

	return &t
}

// SingleDataSourceConfigWithTimeout is the configuration for a single item data
// source that waits for a timeout duration before failing. It represents the
// schema and operations needed to create the single data source.
type SingleDataSourceConfigWithTimeout[TFModel any, APIModel any] struct {
	*BaseSingleDataSourceConfig[TFModel, APIModel]

	// GetAPIModelWithTimeout returns a single API object on which the data
	// source should expose. The function should retry its application logic up
	// until timeout.
	GetAPIModelWithTimeout func(ctx context.Context, tfModel TFModel, client *bastionzero.Client, timeout time.Duration) (*APIModel, error)

	// DefaultTimeout to use if the practitioner does not specify a timeout in
	// the "timeouts" field.
	DefaultTimeout time.Duration
}

type getAPIResult[APIModel any] struct {
	apiObject *APIModel
	err       error
}

// Returns a new single data source given the specified configuration. The
// function panics if the config is invalid. A single data source abstracts
// calling a BastionZero API endpoint that returns a single object.
//
// Reflection is used to add a "timeouts" field to the TF schema which allows
// the practitioner to configure how long to retry calling the BastionZero API
// for a resource.
//
// config.GetAPIModelWithTimeout() should retry its application logic.
// Cancellation occurs if a timeout or interrupt is received.
func NewSingleDataSourceWithTimeout[TFModel any, APIModel any](config *SingleDataSourceConfigWithTimeout[TFModel, APIModel]) datasource.DataSourceWithConfigure {
	if config.RecordSchema == nil {
		panic("RecordSchema cannot be nil")
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
		attributes := config.RecordSchema
		// Add timeouts field
		attributes["timeouts"] = timeouts.Attributes(ctx)

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
		var T TFModel
		var err error

		// Using reflection, add "Timeouts" field
		modelWithTimeouts := dynamicstruct.ExtendStruct(T).
			AddField("Timeouts", timeouts.Value{}, `tfsdk:"timeouts"`).
			Build().
			New()

		// Read Terraform configuration data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, modelWithTimeouts)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Using reflection, copy all fields from modelWithTimeouts, excluding
		// additional "Timeouts" field which does not exist on T, into T
		err = copier.Copy(&T, modelWithTimeouts)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error during conversion to model without timeout",
				fmt.Sprintf("Got: %v during conversion. Please report this issue to the provider developers.", err.Error()),
			)
			return
		}

		// Read() is passed a default timeout to use if no value has been
		// supplied in the Terraform configuration.
		readTimeout, diags := reflect.Indirect(reflect.ValueOf(modelWithTimeouts)).FieldByName("Timeouts").Interface().(timeouts.Value).Read(ctx, config.DefaultTimeout)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Create child context that cancels after readTimeout elapses
		ctx, cancel := context.WithTimeout(ctx, readTimeout)
		defer cancel()

		// Query BastionZero for API object on separate goroutine
		resultCh := make(chan getAPIResult[APIModel], 1)
		go func() {
			defer close(resultCh)
			tflog.Debug(ctx, fmt.Sprintf("Querying for %s", config.PrettyAttributeName))
			apiObject, err := config.GetAPIModelWithTimeout(ctx, T, t.client, readTimeout)
			resultCh <- getAPIResult[APIModel]{apiObject: apiObject, err: err}
		}()

		// Wait for result, timeout, or interrupt
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		var apiObject *APIModel
		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				fmt.Sprintf("Took longer than %s to read", readTimeout),
			)
			return
		case sig := <-sigChan:
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				fmt.Sprintf("Interrupted while reading %s. Caught signal %v", config.PrettyAttributeName, sig),
			)
			return
		case result := <-resultCh:
			apiObject = result.apiObject
			err = result.err
			break
		}

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName))

		// Convert to TFModel
		tfModel, diags := config.FlattenAPIModel(ctx, apiObject, T)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Using reflection, copy all fields from TFModel back into the expected
		// struct stored in the TF state
		err = copier.Copy(modelWithTimeouts, tfModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error during conversion to model with timeout",
				fmt.Sprintf("Got: %v during conversion. Please report this issue to the provider developers.", err.Error()),
			)
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, modelWithTimeouts)...)
	}

	return &t
}
