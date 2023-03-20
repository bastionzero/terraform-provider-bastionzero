package webtarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
)

// webTargetModel maps the web target schema data.
type webTargetModel struct {
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
	RemotePort    types.Number `tfsdk:"remote_port"`
	LocalPort     types.Number `tfsdk:"local_port"`
}

func (t *webTargetModel) SetID(value types.String)              { t.ID = value }
func (t *webTargetModel) SetName(value types.String)            { t.Name = value }
func (t *webTargetModel) SetType(value types.String)            { t.Type = value }
func (t *webTargetModel) SetStatus(value types.String)          { t.Status = value }
func (t *webTargetModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *webTargetModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *webTargetModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *webTargetModel) SetRegion(value types.String)          { t.Region = value }
func (t *webTargetModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

func (t *webTargetModel) SetProxyTargetID(value types.String) { t.ProxyTargetID = value }
func (t *webTargetModel) SetRemoteHost(value types.String)    { t.RemoteHost = value }
func (t *webTargetModel) SetRemotePort(value types.Number)    { t.RemotePort = value }
func (t *webTargetModel) SetLocalPort(value types.Number)     { t.LocalPort = value }

// setWebTargetAttributes populates the TF schema data from a web target API
// object
func setWebTargetAttributes(ctx context.Context, schema *webTargetModel, webTarget *targets.WebTarget) {
	target.SetBaseTargetAttributes(ctx, schema, webTarget)
	target.SetBaseVirtualTargetAttributes(ctx, schema, webTarget)
}

func makeWebTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	webTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Web, opts)
	maps.Copy(webTargetAttributes, target.BaseVirtualTargetDataSourceAttributes(targettype.Web))

	return webTargetAttributes
}
