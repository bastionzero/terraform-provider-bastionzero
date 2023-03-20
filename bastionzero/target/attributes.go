package target

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/targetstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
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
	IsIDComputed bool
	IsIDOptional bool

	IsNameOptional bool
	IsNameComputed bool
}

// BaseTargetDataSourceAttributes returns a map of common TF attributes used by
// the bzero, database, kube, and web data source schemas.
func BaseTargetDataSourceAttributes(targetType targettype.TargetType, opts *BaseTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	baseTargetAttributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    opts.IsIDComputed,
			Required:    opts.IsIDRequired,
			Optional:    opts.IsIDOptional,
			Description: "The target's unique ID.",
		},
		"name": schema.StringAttribute{
			Computed:    opts.IsNameComputed,
			Optional:    opts.IsNameOptional,
			Description: "The target's name.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's type (constant value \"%s\").", targetType),
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
			Description: "The time this target's backing agent last had a transition change in status formatted as a UTC timestamp string in RFC 3339 format. Null if there has not been a single transition change.",
		},
		"agent_version": schema.StringAttribute{
			Computed:    true,
			Description: "The target's backing agent's version.",
		},
		"region": schema.StringAttribute{
			Computed:    true,
			Description: "The BastionZero region that this target has connected to (follows same naming convention as AWS regions).",
		},
		"agent_public_key": schema.StringAttribute{
			Computed:    true,
			Description: "The target's backing agent's public key.",
		},
	}

	return baseTargetAttributes
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
		Description: "Information about the target's backing agent's currently active control channel. Null if the target has no active control channel.",
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
				Description: "The time this control channel connected to the connection node formatted as a UTC timestamp string in RFC 3339 format.",
			},
			"end_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time this control channel disconnected from the connection node formatted as a UTC timestamp string in RFC 3339 format. Null if the control channel is still active.",
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
