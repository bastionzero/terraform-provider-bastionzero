package target

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/targetstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzvalidator"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/typesext"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TargetModelInterface lets you work with common attributes from any kind of
// target model (excluding DAC targets)
type TargetModelInterface interface {
	// SetID sets the target model's ID attribute.
	SetID(value types.String)
	// SetName sets the target model's name attribute.
	SetName(value types.String)
	// SetType sets the target model's type attribute.
	SetType(value types.String)
	// SetStatus sets the target model's status attribute.
	SetStatus(value types.String)
	// SetEnvironmentID sets the target model's environment ID attribute.
	SetEnvironmentID(value types.String)
	// SetLastAgentUpdate sets the target model's last agent update attribute.
	SetLastAgentUpdate(value types.String)
	// SetAgentVersion sets the target model's agent version attribute.
	SetAgentVersion(value types.String)
	// SetRegion sets the target model's region attribute.
	SetRegion(value types.String)
	// SetAgentPublicKey sets the target model's agent public key attribute.
	SetAgentPublicKey(value types.String)
}

// SetBaseTargetAttributes populates base target attributes in the TF schema
// from a base target
func SetBaseTargetAttributes(ctx context.Context, schema TargetModelInterface, baseTarget targets.TargetInterface) {
	schema.SetID(types.StringValue(baseTarget.GetID()))
	schema.SetName(types.StringValue(baseTarget.GetName()))
	schema.SetType(types.StringValue(string(baseTarget.GetTargetType())))
	schema.SetStatus(types.StringValue(string(baseTarget.GetStatus())))
	schema.SetEnvironmentID(types.StringValue(baseTarget.GetEnvironmentID()))

	if baseTarget.GetLastAgentUpdate() != nil {
		schema.SetLastAgentUpdate(types.StringValue(baseTarget.GetLastAgentUpdate().UTC().Format(time.RFC3339)))
	} else {
		schema.SetLastAgentUpdate(types.StringNull())
	}

	schema.SetAgentVersion(types.StringValue(baseTarget.GetAgentVersion()))
	schema.SetRegion(types.StringValue(baseTarget.GetRegion()))
	schema.SetAgentPublicKey(types.StringValue(baseTarget.GetAgentPublicKey()))
}

// BaseTargetDataSourceAttributeOptions are options to use when constructing the
// list of common TF attributes used by the bzero, database, kube, and web data
// source schemas.
type BaseTargetDataSourceAttributeOptions struct {
	IsIDRequired bool
	IsIDOptional bool
	IsIDComputed bool

	IsNameOptional bool
	IsNameComputed bool
}

// BaseTargetDataSourceAttributes returns a map of common TF attributes used by
// the bzero, database, kube, and web data source schemas.
func BaseTargetDataSourceAttributes(targetType targettype.TargetType, opts *BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Required:    opts.IsIDRequired,
			Computed:    opts.IsIDComputed,
			Optional:    opts.IsIDOptional,
			Description: "The target's unique ID.",
			Validators: []validator.String{
				bzvalidator.ValidUUIDV4(),
			},
		},
		"name": schema.StringAttribute{
			Computed:    opts.IsNameComputed,
			Optional:    opts.IsNameOptional,
			Description: "The target's name.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's type (constant value `%s`).", targetType),
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's status %s.", internal.PrettyOneOf(targetstatus.TargetStatusValues())),
			Validators: []validator.String{
				stringvalidator.OneOf(bastionzero.ToStringSlice(targetstatus.TargetStatusValues())...),
			},
		},
		"environment_id": schema.StringAttribute{
			Computed:    true,
			Description: "The target's environment's ID.",
		},
		"last_agent_update": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The time this target's proxy agent last had a transition change in status %s. Null if there has not been a single transition change.", internal.PrettyRFC3339Timestamp()),
		},
		"agent_version": schema.StringAttribute{
			Computed:    true,
			Description: "The target's proxy agent's version.",
		},
		"region": schema.StringAttribute{
			Computed:    true,
			Description: "The BastionZero region that this target has connected to (follows same naming convention as AWS regions).",
		},
		"agent_public_key": schema.StringAttribute{
			Computed:    true,
			Description: "The target's proxy agent's public key.",
		},
	}
}

func TargetDataSourceWithTimeoutMarkdownDescription(baseDescription string, targetType targettype.TargetType) string {
	return fmt.Sprintf("%v"+
		"\n\nSpecify exactly one of `id` or `name`. When specifying a `name`, an error is triggered if more than one %v target is found. "+
		"This data source retries with exponential backoff (provide optional `timeouts.read` [duration](https://pkg.go.dev/time#ParseDuration) to control how long to retry. Defaults to 15 minutes.) until the %v target is found. "+
		"This is useful if there is a chance the target does not exist yet (e.g. the target is in the process of registering to BastionZero).", baseDescription, targetType, targetType)
}

// BaseVirtualTargetDataSourceAttributes returns a map of common TF attributes
// used by the database and web data source schemas.
func BaseVirtualTargetDataSourceAttributes(targetType targettype.TargetType) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"proxy_target_id": schema.StringAttribute{
			Computed:            true,
			Description:         "The target's proxy target's ID (ID of a Bzero or Cluster target).",
			MarkdownDescription: "The target's proxy target's ID (ID of a [Bzero](bzero_target) or [Cluster](cluster_target) target).",
		},
		"remote_host": schema.StringAttribute{
			Computed:    true,
			Description: "The target's hostname or IP address.",
		},
		"remote_port": schema.Int64Attribute{
			Computed:    true,
			Description: fmt.Sprintf("The port of the %v server accessible via the target.", targetType),
		},
		"local_port": schema.Int64Attribute{
			Computed:    true,
			Description: fmt.Sprintf("The port of the %v daemon's localhost server that is spawned on the user's machine on connect. Null if not configured.", targetType),
		},
	}
}

// VirtualTargetModelInterface lets you work with common attributes from any
// kind of virtual target model
type VirtualTargetModelInterface interface {
	// SetProxyTargetID sets the target model's proxy_target_id attribute.
	SetProxyTargetID(value types.String)
	// SetRemoteHost sets the target model's remote_host attribute.
	SetRemoteHost(value types.String)
	// SetRemotePort sets the target model's remote_port attribute.
	SetRemotePort(value types.Int64)
	// SetLocalPort sets the target model's local_port attribute.
	SetLocalPort(value types.Int64)
}

// SetBaseVirtualTargetAttributes populates base virtual target attributes in
// the TF schema from a virtual target
func SetBaseVirtualTargetAttributes(ctx context.Context, schema VirtualTargetModelInterface, virtualTarget targets.VirtualTargetInterface) {
	schema.SetProxyTargetID(types.StringValue(virtualTarget.GetProxyTargetID()))
	schema.SetRemoteHost(types.StringValue(virtualTarget.GetRemoteHost()))
	schema.SetRemotePort(typesext.Int64PointerValue(virtualTarget.GetRemotePort().Value))
	schema.SetLocalPort(typesext.Int64PointerValue(virtualTarget.GetLocalPort().Value))
}

// ControlChannelSummaryModel maps control channel summary data.
type ControlChannelSummaryModel struct {
	ControlChannelID types.String `tfsdk:"control_channel_id"`
	ConnectionNodeID types.String `tfsdk:"connection_node_id"`
	StartTime        types.String `tfsdk:"start_time"`
	EndTime          types.String `tfsdk:"end_time"`
}

func ControlChannelSummaryAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Information about the target's proxy agent's currently active control channel. Null if the target has no active control channel.",
		Attributes: map[string]schema.Attribute{
			"control_channel_id": schema.StringAttribute{
				Computed:    true,
				Description: "The control channel's unique ID.",
			},
			"connection_node_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the connection node that this control channel is connected to.",
			},
			"start_time": schema.StringAttribute{
				Computed:    true,
				Description: fmt.Sprintf("The time this control channel connected to the connection node %s.", internal.PrettyRFC3339Timestamp()),
			},
			"end_time": schema.StringAttribute{
				Computed:    true,
				Description: fmt.Sprintf("The time this control channel disconnected from the connection node %s. Null if the control channel is still active.", internal.PrettyRFC3339Timestamp()),
			},
		},
	}
}

func FlattenControlChannelSummary(ctx context.Context, apiObject *targets.ControlChannelSummary) types.Object {
	attributeTypes, _ := internal.AttributeTypes[ControlChannelSummaryModel](ctx)

	if apiObject == nil {
		return types.ObjectNull(attributeTypes)
	} else {
		var endTime types.String
		if apiObject.EndTime != nil {
			endTime = types.StringValue(apiObject.EndTime.UTC().Format(time.RFC3339))
		} else {
			endTime = types.StringNull()
		}

		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"control_channel_id": types.StringValue(apiObject.ControlChannelID),
			"connection_node_id": types.StringValue(apiObject.ConnectionNodeID),
			"start_time":         types.StringValue(apiObject.StartTime.UTC().Format(time.RFC3339)),
			"end_time":           endTime,
		})
	}
}

type DatabaseAuthenticationConfig struct {
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
				Description: "Cloud service provider hosting the database. Only used for certain types of authentication, such as ServiceAccountInjection.",
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
	attributeTypes, _ := internal.AttributeTypes[ControlChannelSummaryModel](ctx)

	return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"authentication_type":    types.StringPointerValue(apiObject.AuthenticationType),
		"cloud_service_provider": types.StringPointerValue(apiObject.CloudServiceProvider),
		"database":               types.StringPointerValue(apiObject.Database),
		"label":                  types.StringPointerValue(apiObject.Label),
	})
}
