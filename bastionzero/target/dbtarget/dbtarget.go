package dbtarget

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/targetstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/maps"
)

// dbTargetDataSourceModel maps the db target data source schema data.
type dbTargetDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	EnvironmentID   types.String `tfsdk:"environment_id"`
	LastAgentUpdate types.String `tfsdk:"last_agent_update"`
	AgentVersion    types.String `tfsdk:"agent_version"`
	Region          types.String `tfsdk:"region"`
	AgentPublicKey  types.String `tfsdk:"agent_public_key"`

	ProxyTargetID      types.String `tfsdk:"proxy_target_id"`
	ProxyEnvironmentID types.String `tfsdk:"proxy_environment_id"`
	RemoteHost         types.String `tfsdk:"remote_host"`
	RemotePort         types.Int64  `tfsdk:"remote_port"`
	LocalPort          types.Int64  `tfsdk:"local_port"`

	IsSplitCert                  types.Bool   `tfsdk:"is_split_cert"`
	DatabaseType                 types.String `tfsdk:"database_type"`
	DatabaseAuthenticationConfig types.Object `tfsdk:"database_authentication_config"`
}

func (t *dbTargetDataSourceModel) SetID(value types.String)              { t.ID = value }
func (t *dbTargetDataSourceModel) SetName(value types.String)            { t.Name = value }
func (t *dbTargetDataSourceModel) SetType(value types.String)            { t.Type = value }
func (t *dbTargetDataSourceModel) SetStatus(value types.String)          { t.Status = value }
func (t *dbTargetDataSourceModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *dbTargetDataSourceModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *dbTargetDataSourceModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *dbTargetDataSourceModel) SetRegion(value types.String)          { t.Region = value }
func (t *dbTargetDataSourceModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

func (t *dbTargetDataSourceModel) SetProxyTargetID(value types.String) { t.ProxyTargetID = value }
func (t *dbTargetDataSourceModel) SetProxyEnvironmentID(value types.String) {
	t.ProxyEnvironmentID = value
}
func (t *dbTargetDataSourceModel) SetRemoteHost(value types.String) { t.RemoteHost = value }
func (t *dbTargetDataSourceModel) SetRemotePort(value types.Int64)  { t.RemotePort = value }
func (t *dbTargetDataSourceModel) SetLocalPort(value types.Int64)   { t.LocalPort = value }

// setDbTargetDataSourceAttributes populates the TF schema data from a db target
// API object.
func setDbTargetDataSourceAttributes(ctx context.Context, schema *dbTargetDataSourceModel, dbTarget *targets.DatabaseTarget) {
	target.SetBaseTargetAttributes(ctx, schema, dbTarget)
	target.SetBaseVirtualTargetAttributes(ctx, schema, dbTarget)

	schema.IsSplitCert = types.BoolValue(dbTarget.IsSplitCert)
	schema.DatabaseType = types.StringPointerValue(dbTarget.DatabaseType)
	schema.DatabaseAuthenticationConfig = FlattenDatabaseAuthenticationConfig(ctx, &dbTarget.DatabaseAuthenticationConfig)
}

func makeDbTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]datasource_schema.Attribute {
	dbTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Db, opts)
	maps.Copy(dbTargetAttributes, target.BaseVirtualTargetDataSourceAttributes(targettype.Db))
	dbTargetAttributes["database_authentication_config"] = datasource_schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Information about the db target's database authentication configuration.",
		Attributes:  DatabaseAuthenticationConfigAttributes(),
	}
	dbTargetAttributes["is_split_cert"] = datasource_schema.BoolAttribute{
		Computed:           true,
		Description:        "Deprecated. If `true`, this Db target has the split cert feature enabled; `false` otherwise.",
		DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
	}
	dbTargetAttributes["database_type"] = datasource_schema.StringAttribute{
		Computed:           true,
		Description:        "Deprecated. The database's type. Can be null if this Db target does not have the split cert feature enabled (see `is_split_cert`).",
		DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
	}

	return dbTargetAttributes
}

// dbTargetResourceModel maps the db target resource schema data.
type dbTargetResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	EnvironmentID   types.String `tfsdk:"environment_id"`
	LastAgentUpdate types.String `tfsdk:"last_agent_update"`
	AgentVersion    types.String `tfsdk:"agent_version"`
	Region          types.String `tfsdk:"region"`
	AgentPublicKey  types.String `tfsdk:"agent_public_key"`

	ProxyTargetID      types.String `tfsdk:"proxy_target_id"`
	ProxyEnvironmentID types.String `tfsdk:"proxy_environment_id"`
	RemoteHost         types.String `tfsdk:"remote_host"`
	RemotePort         types.Int64  `tfsdk:"remote_port"`
	LocalPort          types.Int64  `tfsdk:"local_port"`

	DatabaseAuthenticationConfig types.Object `tfsdk:"database_authentication_config"`
}

func (t *dbTargetResourceModel) SetID(value types.String)              { t.ID = value }
func (t *dbTargetResourceModel) SetName(value types.String)            { t.Name = value }
func (t *dbTargetResourceModel) SetType(value types.String)            { t.Type = value }
func (t *dbTargetResourceModel) SetStatus(value types.String)          { t.Status = value }
func (t *dbTargetResourceModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *dbTargetResourceModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *dbTargetResourceModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *dbTargetResourceModel) SetRegion(value types.String)          { t.Region = value }
func (t *dbTargetResourceModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

func (t *dbTargetResourceModel) SetProxyTargetID(value types.String) { t.ProxyTargetID = value }
func (t *dbTargetResourceModel) SetProxyEnvironmentID(value types.String) {
	t.ProxyEnvironmentID = value
}
func (t *dbTargetResourceModel) SetRemoteHost(value types.String) { t.RemoteHost = value }
func (t *dbTargetResourceModel) SetRemotePort(value types.Int64)  { t.RemotePort = value }
func (t *dbTargetResourceModel) SetLocalPort(value types.Int64)   { t.LocalPort = value }

// setDbTargetResourceAttributes populates the TF schema data from a db target
// API object.
func setDbTargetResourceAttributes(ctx context.Context, schema *dbTargetResourceModel, dbTarget *targets.DatabaseTarget) {
	target.SetBaseTargetAttributes(ctx, schema, dbTarget)
	target.SetBaseVirtualTargetAttributes(ctx, schema, dbTarget)

	schema.DatabaseAuthenticationConfig = FlattenDatabaseAuthenticationConfig(ctx, &dbTarget.DatabaseAuthenticationConfig)
}

func makeDbTargetResourceSchema(ctx context.Context) map[string]resource_schema.Attribute {
	// Valid constants for fields in `database_authentication_config`. These are
	// the only values the BastionZero backend currently accepts
	validAuthenticationTypes := []string{
		dbauthconfig.Default,
		dbauthconfig.SplitCert,
		dbauthconfig.ServiceAccountInjection,
	}
	validCloudServiceProviders := []string{
		dbauthconfig.AWS,
		dbauthconfig.GCP,
	}
	validDatabases := []string{
		dbauthconfig.CockroachDB,
		dbauthconfig.MicrosoftSQLServer,
		dbauthconfig.MongoDB,
		dbauthconfig.MySQL,
		dbauthconfig.Postgres,
	}

	return map[string]resource_schema.Attribute{
		"id": resource_schema.StringAttribute{
			Computed:    true,
			Description: "The target's unique ID.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": resource_schema.StringAttribute{
			Required:    true,
			Description: "The target's name.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"type": resource_schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's type (constant value `%s`).", targettype.Db),
			Default:     stringdefault.StaticString(string(targettype.Db)),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"status": resource_schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's status %s.", internal.PrettyOneOf(targetstatus.TargetStatusValues())),
			Validators: []validator.String{
				stringvalidator.OneOf(bastionzero.ToStringSlice(targetstatus.TargetStatusValues())...),
			},
		},
		"environment_id": resource_schema.StringAttribute{
			Required:    true,
			Description: "The target's environment's ID.",
			Validators: []validator.String{
				bzvalidator.ValidUUIDV4(),
			},
		},
		"last_agent_update": resource_schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The time this target's proxy agent last had a transition change in status %s. Null if there has not been a single transition change.", internal.PrettyRFC3339Timestamp()),
		},
		"agent_version": resource_schema.StringAttribute{
			Computed:    true,
			Description: "The target's proxy agent's version.",
		},
		"region": resource_schema.StringAttribute{
			Computed:    true,
			Description: "The BastionZero region that this target has connected to (follows same naming convention as AWS regions).",
		},
		"agent_public_key": resource_schema.StringAttribute{
			Computed:    true,
			Description: "The target's proxy agent's public key.",
		},
		"proxy_target_id": resource_schema.StringAttribute{
			Optional:            true,
			Description:         "The target's proxy target's ID (ID of a Bzero or Cluster target).",
			MarkdownDescription: "The target's proxy target's ID (ID of a [Bzero](bzero_target) or [Cluster](cluster_target) target).",
			Validators: []validator.String{
				bzvalidator.ValidUUIDV4(),
			},
		},
		"proxy_environment_id": resource_schema.StringAttribute{
			Optional:    true,
			Description: "The target's proxy environment's ID (ID of the backing proxy environment).",
			Validators: []validator.String{
				bzvalidator.ValidUUIDV4(),
			},
		},
		"remote_host": resource_schema.StringAttribute{
			Required:    true,
			Description: "The target's hostname or IP address.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"remote_port": resource_schema.Int64Attribute{
			Required:    true,
			Description: fmt.Sprintf("The port of the %v server accessible via the target. This field is required for all databases; however, if `database_authentication_config.cloud_service_provider` is equal to `%v`, then the value will be ignored when connecting to the database.", targettype.Db, dbauthconfig.GCP),
		},
		"local_port": resource_schema.Int64Attribute{
			Optional:    true,
			Description: fmt.Sprintf("The port of the %v daemon's localhost server that is spawned on the user's machine on connect. If this attribute is left unconfigured, an available port will be chosen when the target is connected to.", targettype.Db),
		},
		"database_authentication_config": resource_schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			// Important to set this default, even though backend sets it for us
			// when unspecified, because it makes it easier to keep the TF state
			// and TF config in sync with the backend. Removing this code will
			// make it so terraform thinks any TF config that does not set
			// db_auth_config require an update to the remote state when it
			// really does not
			Default: objectdefault.StaticValue(FlattenDatabaseAuthenticationConfig(ctx,
				&dbauthconfig.DatabaseAuthenticationConfig{
					AuthenticationType: bastionzero.PtrTo(dbauthconfig.Default),
					Label:              bastionzero.PtrTo("None"),
				}),
			),
			Description: "Information about the db target's database authentication configuration. If this attribute is left unconfigured, the target is configured with the default, non-passwordless database configuration.",
			Attributes: map[string]resource_schema.Attribute{
				"authentication_type": resource_schema.StringAttribute{
					Required:    true,
					Description: fmt.Sprintf("The type of authentication used when connecting to the database %s.", internal.PrettyOneOf(validAuthenticationTypes)),
					Validators: []validator.String{
						stringvalidator.OneOf(validAuthenticationTypes...),
					},
				},
				"cloud_service_provider": resource_schema.StringAttribute{
					Optional:    true,
					Description: fmt.Sprintf("Cloud service provider hosting the database %s. Only used for certain types of authentication (`authentication_type`), such as `ServiceAccountInjection`.", internal.PrettyOneOf(validCloudServiceProviders)),
					Validators: []validator.String{
						stringvalidator.OneOf(validCloudServiceProviders...),
					},
				},
				"database": resource_schema.StringAttribute{
					Optional:    true,
					Description: fmt.Sprintf("The type of database running on the target %s.", internal.PrettyOneOf(validDatabases)),
					Validators: []validator.String{
						stringvalidator.OneOf(validDatabases...),
					},
				},
				"label": resource_schema.StringAttribute{
					Optional:    true,
					Description: "User-friendly label for this database authentication configuration.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
	}
}

type DatabaseAuthenticationConfigModel struct {
	AuthenticationType   types.String `tfsdk:"authentication_type"`
	CloudServiceProvider types.String `tfsdk:"cloud_service_provider"`
	Database             types.String `tfsdk:"database"`
	Label                types.String `tfsdk:"label"`
}

func DatabaseAuthenticationConfigAttributes() map[string]datasource_schema.Attribute {
	return map[string]datasource_schema.Attribute{
		"authentication_type": datasource_schema.StringAttribute{
			Computed:    true,
			Description: "The type of authentication used when connecting to the database.",
		},
		"cloud_service_provider": datasource_schema.StringAttribute{
			Computed:    true,
			Description: "Cloud service provider hosting the database. Only used for certain types of authentication, such as `ServiceAccountInjection`.",
		},
		"database": datasource_schema.StringAttribute{
			Computed:    true,
			Description: "The type of database running on the target.",
		},
		"label": datasource_schema.StringAttribute{
			Computed:    true,
			Description: "User-friendly label for this database authentication configuration.",
		},
	}
}

func ExpandDatabaseAuthenticationConfig(ctx context.Context, tfObject types.Object) *dbauthconfig.DatabaseAuthenticationConfig {
	return internal.ExpandFrameworkObject(ctx, tfObject, func(m DatabaseAuthenticationConfigModel) *dbauthconfig.DatabaseAuthenticationConfig {
		return &dbauthconfig.DatabaseAuthenticationConfig{
			AuthenticationType:   m.AuthenticationType.ValueStringPointer(),
			CloudServiceProvider: m.CloudServiceProvider.ValueStringPointer(),
			Database:             m.Database.ValueStringPointer(),
			Label:                m.Label.ValueStringPointer(),
		}
	})
}

func FlattenDatabaseAuthenticationConfig(ctx context.Context, apiObject *dbauthconfig.DatabaseAuthenticationConfig) types.Object {
	attributeTypes, _ := internal.AttributeTypes[DatabaseAuthenticationConfigModel](ctx)

	return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"authentication_type":    types.StringPointerValue(apiObject.AuthenticationType),
		"cloud_service_provider": types.StringPointerValue(apiObject.CloudServiceProvider),
		"database":               types.StringPointerValue(apiObject.Database),
		"label":                  types.StringPointerValue(apiObject.Label),
	})
}

func readDbTarget(ctx context.Context, schema *dbTargetResourceModel, client *bastionzero.Client) (found bool, diags diag.Diagnostics) {
	if schema.ID.IsUnknown() || schema.ID.IsNull() {
		diags.AddError(
			"Unexpected null ID in schema",
			"Expected ID to be set. Please report this issue to the provider developers.",
		)
		return false, diags
	}

	// Get refreshed db target value from BastionZero
	tflog.Debug(ctx, "Querying for db target")
	dbTarget, _, err := client.Targets.GetDatabaseTarget(ctx, schema.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		return false, diags
	} else if err != nil {
		diags.AddError(
			"Error reading db target",
			"Could not read db target, unexpected error: "+err.Error())
		return false, diags
	}
	tflog.Debug(ctx, "Queried for db target")

	setDbTargetResourceAttributes(ctx, schema, dbTarget)
	return true, diags
}
