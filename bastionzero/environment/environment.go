package environment

import (
	"context"
	"net/http"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	bzplanmodifier "github.com/bastionzero/terraform-provider-bastionzero/internal/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// environmentModel maps the environment schema data.
type environmentModel struct {
	ID                         types.String `tfsdk:"id"`
	OrganizationID             types.String `tfsdk:"organization_id"`
	IsDefault                  types.Bool   `tfsdk:"is_default"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	TimeCreated                types.String `tfsdk:"time_created"`
	OfflineCleanupTimeoutHours types.Int64  `tfsdk:"offline_cleanup_timeout_hours"`
	Targets                    types.Map    `tfsdk:"targets"` // key is target id. value is environmentTargetModel
}

// environmentTargetModel maps target summary data.
type environmentTargetModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

var (
	environmentTargetModelAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
)

// setEnvironmentAttributes populates the TF schema data from an environment
func setEnvironmentAttributes(ctx context.Context, schema *environmentModel, env *environments.Environment) diag.Diagnostics {
	schema.Name = types.StringValue(env.Name)

	// Use StringEmptyIsNullValue to normalize "" to Terraform Null value (since
	// the schema says description is optional/nullable)
	// schema.Description = typesext.StringEmptyIsNullValue(&env.Description)

	// Preserve null in TF schema. We say that "" is semantically equivalent to
	// null for the environment schema
	if schema.Description.IsNull() && env.Description == "" {
		schema.Description = types.StringNull()
	} else {
		schema.Description = types.StringValue(env.Description)
	}

	schema.OfflineCleanupTimeoutHours = types.Int64Value(int64(env.OfflineCleanupTimeoutHours))

	schema.ID = types.StringValue(env.ID)
	schema.OrganizationID = types.StringValue(env.OrganizationID)
	schema.IsDefault = types.BoolValue(env.IsDefault)
	schema.TimeCreated = types.StringValue(env.TimeCreated.UTC().Format(time.RFC3339))

	targetsMap := make(map[string]environmentTargetModel)
	for _, target := range env.Targets {
		targetsMap[target.ID] = environmentTargetModel{
			ID:   types.StringValue(target.ID),
			Type: types.StringValue(string(target.Type)),
		}
	}
	targets, diags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: environmentTargetModelAttrTypes}, targetsMap)
	schema.Targets = targets

	return diags
}

func makeEnvironmentResourceSchema() map[string]schema.Attribute {
	res := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				// An environment's ID remains the same after an update is made
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The environment's unique ID.",
		},
		"organization_id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The environment's organization's ID.",
		},
		"is_default": schema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
			Description: "If true, this environment is the default environment. False otherwise.",
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "The environment's name.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Description: "The environment's description.",
		},
		"time_created": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "The time this environment was created in BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
		},
		"offline_cleanup_timeout_hours": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The amount of time (in hours) to wait until offline targets are automatically removed by BastionZero (Defaults to 90 days).",
			PlanModifiers: []planmodifier.Int64{
				// Default to 90 days like in webapp
				bzplanmodifier.Int64DefaultValue(types.Int64Value(24 * 90)),
			},
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"targets": schema.MapNestedAttribute{
			Description: "Map of targets that belong to this environment. The map is keyed by a target's unique ID.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "The target's unique ID.",
					},
					"type": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: `This value is one of "Bzero", "Cluster", "DynamicAccessConfig", "Web", or "Db"`,
					},
				},
			},
		},
	}

	return res
}

func readEnvironment(ctx context.Context, schema *environmentModel, client *bastionzero.Client) (found bool, diags diag.Diagnostics) {
	if schema.ID.IsUnknown() || schema.ID.IsNull() {
		diags.AddError(
			"Unexpected null ID in schema",
			"Expected ID to be set. Please report this issue to the provider developers.",
		)
		return false, diags
	}

	// Get refreshed environment value from BastionZero
	tflog.Debug(ctx, "Querying for environment")
	env, httpResp, err := client.Environments.GetEnvironment(ctx, schema.ID.ValueString())
	if httpResp.StatusCode == http.StatusNotFound {
		return false, diags
	}
	if err != nil {
		diags.AddError(
			"Error reading environment",
			"Could not read environment, unexpected error: "+err.Error())
		return false, diags
	}
	tflog.Debug(ctx, "Queried for environment")

	diags.Append(setEnvironmentAttributes(ctx, schema, env)...)
	return true, diags
}
