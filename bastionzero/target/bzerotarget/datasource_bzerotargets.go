package bzerotarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
)

// bzeroTargetModel maps the bzero target schema data.
type bzeroTargetModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
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
func (t *bzeroTargetModel) SetStatus(value types.String)          { t.Status = value }
func (t *bzeroTargetModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *bzeroTargetModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *bzeroTargetModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *bzeroTargetModel) SetRegion(value types.String)          { t.Region = value }
func (t *bzeroTargetModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

func NewBzeroTargetsDataSource() datasource.DataSource {
	bzeroTargetAttributes := map[string]schema.Attribute{
		"control_channel": target.ControlChannelSummaryAttribute(),
	}
	// Add common base target attributes
	maps.Copy(bzeroTargetAttributes, target.BaseTargetDataSourceAttributes())

	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[bzeroTargetModel, targets.BzeroTarget]{
		RecordSchema:        bzeroTargetAttributes,
		ResultAttributeName: "bzero_targets",
		PrettyAttributeName: "Bzero targets",
		FlattenAPIModel: func(ctx context.Context, apiObject targets.BzeroTarget) (state *bzeroTargetModel, diags diag.Diagnostics) {
			state = new(bzeroTargetModel)
			target.SetBaseTargetAttributes(ctx, state, &apiObject)
			state.ControlChannel = target.FlattenControlChannelSummary(ctx, apiObject.ControlChannel)
			return
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.BzeroTarget, error) {
			targets, _, err := client.Targets.ListBzeroTargets(ctx)
			return targets, err
		},
		Description: "Get a list of all Bzero targets in your BastionZero organization.",
	})
}
