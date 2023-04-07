package clustertarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clusterTargetModel maps the cluster target schema data.
type clusterTargetModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Type              types.String `tfsdk:"type"`
	Status            types.String `tfsdk:"status"`
	EnvironmentID     types.String `tfsdk:"environment_id"`
	LastAgentUpdate   types.String `tfsdk:"last_agent_update"`
	AgentVersion      types.String `tfsdk:"agent_version"`
	Region            types.String `tfsdk:"region"`
	AgentPublicKey    types.String `tfsdk:"agent_public_key"`
	ControlChannel    types.Object `tfsdk:"control_channel"`
	ValidClusterUsers types.Set    `tfsdk:"valid_cluster_users"`
}

func (t *clusterTargetModel) SetID(value types.String)              { t.ID = value }
func (t *clusterTargetModel) SetName(value types.String)            { t.Name = value }
func (t *clusterTargetModel) SetType(value types.String)            { t.Type = value }
func (t *clusterTargetModel) SetStatus(value types.String)          { t.Status = value }
func (t *clusterTargetModel) SetEnvironmentID(value types.String)   { t.EnvironmentID = value }
func (t *clusterTargetModel) SetLastAgentUpdate(value types.String) { t.LastAgentUpdate = value }
func (t *clusterTargetModel) SetAgentVersion(value types.String)    { t.AgentVersion = value }
func (t *clusterTargetModel) SetRegion(value types.String)          { t.Region = value }
func (t *clusterTargetModel) SetAgentPublicKey(value types.String)  { t.AgentPublicKey = value }

// setClusterTargetAttributes populates the TF schema data from a cluster target
// API object.
func setClusterTargetAttributes(ctx context.Context, schema *clusterTargetModel, clusterTarget *targets.ClusterTarget) {
	target.SetBaseTargetAttributes(ctx, schema, clusterTarget)
	schema.ControlChannel = target.FlattenControlChannelSummary(ctx, clusterTarget.ControlChannel)
	schema.ValidClusterUsers = internal.FlattenFrameworkSet(ctx, types.StringType, clusterTarget.ValidClusterUsers, func(user string) attr.Value { return types.StringValue(user) })
}

func makeClusterTargetDataSourceSchema(opts *target.BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	clusterTargetAttributes := target.BaseTargetDataSourceAttributes(targettype.Bzero, opts)
	clusterTargetAttributes["control_channel"] = target.ControlChannelSummaryAttribute()
	clusterTargetAttributes["valid_cluster_users"] = schema.SetAttribute{
		Computed:    true,
		Description: "Set of Kubernetes user [subjects](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-subjects) that have been extracted from RoleBindings or ClusterRoleBindings defined in the cluster.",
		ElementType: types.StringType,
	}

	return clusterTargetAttributes
}
