package dbtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
)

// dbTargetModel maps the db target schema data.
type dbTargetModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	EnvironmentID   types.String `tfsdk:"environment_id"`
	LastAgentUpdate types.String `tfsdk:"last_agent_update"`
	AgentVersion    types.String `tfsdk:"agent_version"`
	Region          types.String `tfsdk:"region"`
	AgentPublicKey  types.String `tfsdk:"agent_public_key"`

	ProxyTargetID                types.String `tfsdk:"proxy_target_id"`
	ProxyEnvironmentID           types.String `tfsdk:"proxy_environment_id"`
	RemoteHost                   types.String `tfsdk:"remote_host"`
	RemotePort                   types.Int64  `tfsdk:"remote_port"`
	LocalPort                    types.Int64  `tfsdk:"local_port"`
	IsSplitCert                  types.Bool   `tfsdk:"is_split_cert"`
	DatabaseType                 types.String `tfsdk:"database_type"`
	DatabaseAuthenticationConfig types.Object `tfsdk:"database_authentication_config"`
}

func (t *dbTargetModel) SetID(value types.String)              { t.ID = value }
func (t *dbTargetModel) SetName(value types.String)            { t.Name = value }
func (t *dbTargetModel) SetType(value types.String)            { t.Type = value }
func (t *dbTargetModel) SetStatus(value types.String)          { t.Status = value }
func (t *dbTargetModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *dbTargetModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *dbTargetModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *dbTargetModel) SetRegion(value types.String)          { t.Region = value }
func (t *dbTargetModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

func (t *dbTargetModel) SetProxyTargetID(value types.String)      { t.ProxyTargetID = value }
func (t *dbTargetModel) SetProxyEnvironmentID(value types.String) { t.ProxyEnvironmentID = value }
func (t *dbTargetModel) SetRemoteHost(value types.String)         { t.RemoteHost = value }
func (t *dbTargetModel) SetRemotePort(value types.Int64)          { t.RemotePort = value }
func (t *dbTargetModel) SetLocalPort(value types.Int64)           { t.LocalPort = value }

// setDbTargetAttributes populates the TF schema data from a db target API
// object.
func setDbTargetAttributes(ctx context.Context, schema *dbTargetModel, dbTarget *targets.DatabaseTarget) {
	target.SetBaseTargetAttributes(ctx, schema, dbTarget)
	target.SetBaseVirtualTargetAttributes(ctx, schema, dbTarget)

	schema.IsSplitCert = types.BoolValue(dbTarget.IsSplitCert)
	schema.DatabaseType = types.StringPointerValue(dbTarget.DatabaseType)
	schema.DatabaseAuthenticationConfig = FlattenDatabaseAuthenticationConfig(ctx, dbTarget.DatabaseAuthenticationConfig)
}

func makeDbTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	dbTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Db, opts)
	maps.Copy(dbTargetAttributes, target.BaseVirtualTargetDataSourceAttributes(targettype.Db))
	dbTargetAttributes["database_authentication_config"] = DatabaseAuthenticationConfigAttribute()
	dbTargetAttributes["is_split_cert"] = schema.BoolAttribute{
		Computed:           true,
		Description:        "Deprecated. If `true`, this Db target has the split cert feature enabled; `false` otherwise.",
		DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
	}
	dbTargetAttributes["database_type"] = schema.StringAttribute{
		Computed:           true,
		Description:        "Deprecated. The database's type. Can be null if this Db target does not have the split cert feature enabled (see `is_split_cert`).",
		DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
	}

	return dbTargetAttributes
}

type DatabaseAuthenticationConfigModel struct {
	AuthenticationType   types.String `tfsdk:"authentication_type"`
	CloudServiceProvider types.String `tfsdk:"cloud_service_provider"`
	Database             types.String `tfsdk:"database"`
	Label                types.String `tfsdk:"label"`
}

func DatabaseAuthenticationConfigAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Information about the db target's database authentication configuration.",
		Attributes: map[string]schema.Attribute{
			"authentication_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of authentication used when connecting to the database.",
			},
			"cloud_service_provider": schema.StringAttribute{
				Computed:    true,
				Description: "Cloud service provider hosting the database. Only used for certain types of authentication, such as `ServiceAccountInjection`.",
			},
			"database": schema.StringAttribute{
				Computed:    true,
				Description: "The type of database running on the target.",
			},
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "User-friendly label for this database authentication configuration.",
			},
		},
	}
}

func FlattenDatabaseAuthenticationConfig(ctx context.Context, apiObject dbauthconfig.DatabaseAuthenticationConfig) types.Object {
	attributeTypes, _ := internal.AttributeTypes[DatabaseAuthenticationConfigModel](ctx)

	return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"authentication_type":    types.StringPointerValue(apiObject.AuthenticationType),
		"cloud_service_provider": types.StringPointerValue(apiObject.CloudServiceProvider),
		"database":               types.StringPointerValue(apiObject.Database),
		"label":                  types.StringPointerValue(apiObject.Label),
	})
}
