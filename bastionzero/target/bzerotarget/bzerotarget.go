package bzerotarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bzeroTargetModel maps the bzero target schema data.
type bzeroTargetModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Status          types.String `tfsdk:"status"`
	EnvironmentID   types.String `tfsdk:"environment_id"`
	LastAgentUpdate types.String `tfsdk:"last_agent_update"`
	AgentVersion    types.String `tfsdk:"agent_version"`
	Region          types.String `tfsdk:"region"`
	AgentPublicKey  types.String `tfsdk:"agent_public_key"`
	ControlChannel  types.Object `tfsdk:"control_channel"`
}

func (t *bzeroTargetModel) SetID(value types.String)              { t.ID = value }
func (t *bzeroTargetModel) SetName(value types.String)            { t.Name = value }
func (t *bzeroTargetModel) SetType(value types.String)            { t.Type = value }
func (t *bzeroTargetModel) SetStatus(value types.String)          { t.Status = value }
func (t *bzeroTargetModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *bzeroTargetModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *bzeroTargetModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *bzeroTargetModel) SetRegion(value types.String)          { t.Region = value }
func (t *bzeroTargetModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

// setBzeroTargetAttributes populates the TF schema data from a bzero target API
// object.
func setBzeroTargetAttributes(ctx context.Context, schema *bzeroTargetModel, bzeroTarget *targets.BzeroTarget) {
	target.SetBaseTargetAttributes(ctx, schema, bzeroTarget)
	schema.ControlChannel = target.FlattenControlChannelSummary(ctx, bzeroTarget.ControlChannel)
}

func makeBzeroTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	bzeroTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Bzero, opts)
	bzeroTargetAttributes["control_channel"] = target.ControlChannelSummaryAttribute()

	return bzeroTargetAttributes
}
