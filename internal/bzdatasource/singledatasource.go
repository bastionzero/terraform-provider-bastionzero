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
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jinzhu/copier"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

// TFSingleDataSourceModel is a struct that models a collection of TF schema
// attributes that can be a mix of Computed, Required, and Optional attributes.
type TFSingleDataSourceModel = interface{}

type BaseSingleDataSourceConfig[T TFSingleDataSourceModel, T2 APIModel] struct {
	// RecordSchema is the TF schema that models a single instance of the API
	// object. There should be a key for each field defined in
	// TFSingleDataSourceModel. Required.
	RecordSchema map[string]schema.Attribute

	// MetadataTypeName is the suffix to use for the name of the data source.
	// Cannot be the empty string.
	MetadataTypeName string

	// PrettyAttributeName is the name of the attribute used for logging and
	// documentation purposes. Cannot be the empty string.
	PrettyAttributeName string

	// FlattenAPIModel takes a model returned from the GetAPIModel function and
	// uses it to update the TF model.
	FlattenAPIModel func(ctx context.Context, apiObject *T2, tfModel *T) diag.Diagnostics

	// GetAPIModel returns a single API object on which the data source should
	// expose.
	GetAPIModel func(ctx context.Context, tfModel T, client *bastionzero.Client) (*T2, error)

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

// Validate checks for errors in the base single data source config
func (c *BaseSingleDataSourceConfig[T, T2]) Validate() error {
	if c.RecordSchema == nil {
		return fmt.Errorf("RecordSchema cannot be nil")
	}
	if c.MetadataTypeName == "" {
		return fmt.Errorf("MetadataTypeName cannot be empty")
	}
	if c.PrettyAttributeName == "" {
		return fmt.Errorf("PrettyAttributeName cannot be empty")
	}

	return nil
}

// SingleDataSource is a data source that calls a BastionZero API endpoint which
// returns a single object. It abstracts common, boilerplate code that typically
// accompanies a data source that exposes a single item. It is assumed the TF
// model contains some identification attribute(s) and/or parameters that can be
// used when querying the BastionZero API.
type SingleDataSource datasource.DataSourceWithConfigure

// SingleDataSourceConfig is the configuration for a single data source. It
// represents the schema and operations needed to create the data source.
type SingleDataSourceConfig[T TFSingleDataSourceModel, T2 APIModel] struct {
	*BaseSingleDataSourceConfig[T, T2]
}

// NewSingleDataSource creates a SingleDataSource. The function panics if the
// config is invalid.
func NewSingleDataSource[T TFSingleDataSourceModel, T2 APIModel](config *SingleDataSourceConfig[T, T2]) SingleDataSource {
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
		var T T
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
		diags := config.FlattenAPIModel(ctx, apiObject, &T)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &T)...)
	}

	return &t
}

// SingleDataSourceWithTimeout is a data source with operational semantics
// similar to bzdatasource.SingleDataSource. It additionally abstracts calling a
// BastionZero API endpoint many times until it returns a valid response to
// expose through the data source.
//
// GetAPIModel() is called with exponential backoff until a result is returned,
// or a timeout occurs. The timeout can be configured by the practitioner. If
// the error is fatal and should short circuit, have GetAPIModel() return an
// error of type backoff.PermanentError.
type SingleDataSourceWithTimeout datasource.DataSourceWithConfigure

// SingleDataSourceWithTimeoutConfig is the configuration for a single data
// source with timeout. It represents the schema and operations needed to create
// the data source.
type SingleDataSourceWithTimeoutConfig[T TFSingleDataSourceModel, T2 APIModel] struct {
	*BaseSingleDataSourceConfig[T, T2]

	// DefaultTimeout to use if the practitioner does not specify a timeout in
	// the "timeouts" field.
	DefaultTimeout time.Duration
}

// NewSingleDataSourceWithTimeout creates a SingleDataSourceWithTimeout. The
// function panics if the config is invalid.
func NewSingleDataSourceWithTimeout[T TFSingleDataSourceModel, T2 APIModel](config *SingleDataSourceWithTimeoutConfig[T, T2]) SingleDataSourceWithTimeout {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	t := struct{ protoDataSource }{}
	t.metadataFunc = func(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
		resp.TypeName = req.ProviderTypeName + fmt.Sprintf("_%s", config.MetadataTypeName)
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
		var model T
		var err error

		// Using reflection, add "Timeouts" field
		modelWithTimeouts := dynamicstruct.ExtendStruct(model).
			AddField("Timeouts", timeouts.Value{}, `tfsdk:"timeouts"`).
			Build().
			New()

		// Read Terraform configuration data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, modelWithTimeouts)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Using reflection, copy all fields from modelWithTimeouts, excluding
		// additional "Timeouts" field which does not exist on the model, into
		// the model
		err = copier.Copy(&model, modelWithTimeouts)
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

		// Create linked child context that we can cancel under our own
		// conditions in addition to the Terraform framework's context.
		childCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Spawn goroutine that listens for interrupts from user
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			select {
			case <-sigChan:
				cancel()
				return
			case <-childCtx.Done():
				// Must return to not leak this goroutine (in case no interrupt
				// received)
				return
			}
		}()

		// Perform API call with backoff
		backOffConfig := backoff.NewExponentialBackOff()
		// Stop trying after timeout is hit
		backOffConfig.MaxElapsedTime = readTimeout

		apiObject, err := backoff.RetryNotifyWithData(
			func() (*T2, error) {
				apiObject, err := config.GetAPIModel(childCtx, model, t.client)
				return apiObject, err
			},
			// Init backoff config with child context, so that we can cancel it
			// due to interrupt
			backoff.WithContext(backOffConfig, childCtx),
			// Log message
			func(err error, dur time.Duration) {
				tflog.Info(ctx, fmt.Sprintf("%v. Retrying in %s...", err, dur))
			},
		)

		// Error from API server, timeout, or context cancelled
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error reading %s", config.PrettyAttributeName),
				err.Error(),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Queried for %s", config.PrettyAttributeName))

		// Convert to TFModel
		diags = config.FlattenAPIModel(ctx, apiObject, &model)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Using reflection, copy all fields from TFModel back into the expected
		// struct stored in the TF state
		err = copier.Copy(modelWithTimeouts, &model)
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
