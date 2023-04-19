package dbtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
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

	ProxyTargetID types.String `tfsdk:"proxy_target_id"`
	RemoteHost    types.String `tfsdk:"remote_host"`
	RemotePort    types.Int64  `tfsdk:"remote_port"`
	LocalPort     types.Int64  `tfsdk:"local_port"`
	IsSplitCert   types.Bool   `tfsdk:"is_split_cert"`
	DatabaseType  types.String `tfsdk:"database_type"`
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

func (t *dbTargetModel) SetProxyTargetID(value types.String) { t.ProxyTargetID = value }
func (t *dbTargetModel) SetRemoteHost(value types.String)    { t.RemoteHost = value }
func (t *dbTargetModel) SetRemotePort(value types.Int64)     { t.RemotePort = value }
func (t *dbTargetModel) SetLocalPort(value types.Int64)      { t.LocalPort = value }

// setDbTargetAttributes populates the TF schema data from a db target API
// object.
func setDbTargetAttributes(ctx context.Context, schema *dbTargetModel, dbTarget *targets.DatabaseTarget) {
	target.SetBaseTargetAttributes(ctx, schema, dbTarget)
	target.SetBaseVirtualTargetAttributes(ctx, schema, dbTarget)

	schema.IsSplitCert = types.BoolValue(dbTarget.IsSplitCert)
	schema.DatabaseType = types.StringPointerValue(dbTarget.DatabaseType)
}

func makeDbTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	dbTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Db, opts)
	maps.Copy(dbTargetAttributes, target.BaseVirtualTargetDataSourceAttributes(targettype.Db))
	dbTargetAttributes["is_split_cert"] = schema.BoolAttribute{
		Computed:    true,
		Description: "If true, this Db target has the split cert feature enabled. False otherwise.",
	}
	dbTargetAttributes["database_type"] = schema.StringAttribute{
		Computed:    true,
		Description: "The database's type. Can be null if this Db target does not have the split cert feature enabled (see `is_split_cert`).",
	}

	return dbTargetAttributes
}
