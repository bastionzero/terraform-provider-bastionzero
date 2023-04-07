package dactarget

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dacstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dacTargetModel maps the DAC target schema data.
type dacTargetModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	StartWebhook  types.String `tfsdk:"start_webhook"`
	StopWebhook   types.String `tfsdk:"stop_webhook"`
	HealthWebhook types.String `tfsdk:"health_webhook"`
	Status        types.String `tfsdk:"status"`
}

// setDacTargetAttributes populates the TF schema data from a web target API
// object
func setDacTargetAttributes(ctx context.Context, schema *dacTargetModel, dacTarget *targets.DynamicAccessConfiguration) {
	schema.ID = types.StringValue(dacTarget.ID)
	schema.Name = types.StringValue(dacTarget.Name)
	schema.Type = types.StringValue(string(targettype.DynamicAccessConfig))
	schema.EnvironmentID = types.StringValue(dacTarget.EnvironmentId)
	schema.StartWebhook = types.StringValue(dacTarget.StartWebhook)
	schema.StopWebhook = types.StringValue(dacTarget.StopWebhook)
	schema.HealthWebhook = types.StringValue(dacTarget.HealthWebhook)
	schema.Status = types.StringValue(string(dacTarget.Status))
}

type dacTargetDataSourceAttributeOptions struct {
	IsIDComputed bool
	IsIDRequired bool
}

func makeDacTargetDataSourceSchema(opts *dacTargetDataSourceAttributeOptions) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    opts.IsIDComputed,
			Required:    opts.IsIDRequired,
			Description: "The DAC's unique ID.",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "The DAC's name.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The target's type (constant value `%s`).", targettype.DynamicAccessConfig),
		},
		"environment_id": schema.StringAttribute{
			Computed:    true,
			Description: "The DAC's environment's ID.",
		},
		"start_webhook": schema.StringAttribute{
			Computed:    true,
			Description: "URL for the dynamic access provisioning server's webhook that starts a new instance.",
		},
		"stop_webhook": schema.StringAttribute{
			Computed:    true,
			Description: "URL for the dynamic access provisioning server's webhook that stops a new instance.",
		},
		"health_webhook": schema.StringAttribute{
			Computed:    true,
			Description: "URL for the dynamic access provisioning server's webhook that does a health check.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The DAC's status %s.", internal.PrettyOneOf(dacstatus.DynamicAccessConfigurationStatusValues())),
			Validators: []validator.String{
				stringvalidator.OneOf(bastionzero.ToStringSlice(dacstatus.DynamicAccessConfigurationStatusValues())...),
			},
		},
	}
}
