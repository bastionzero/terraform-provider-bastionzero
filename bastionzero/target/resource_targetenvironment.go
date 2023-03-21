package target

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &targetEnvironmentResource{}
	_ resource.ResourceWithConfigure = &targetEnvironmentResource{}
)

func NewTargetEnvironmentResource() resource.Resource {
	return &targetEnvironmentResource{}
}

// targetEnvironmentModel maps the target environment schema data.
type targetEnvironmentModel struct {
	TargetID      types.String `tfsdk:"target_id"`
	TargetType    types.String `tfsdk:"target_type"`
	EnvironmentID types.String `tfsdk:"environment_id"`
}

// targetEnvironmentResource is the resource implementation.
type targetEnvironmentResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *targetEnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*bastionzero.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource configure type",
			fmt.Sprintf("Expected *bastionzero.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the target connect policy resource type name.
func (r *targetEnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target_environment"
}

// Schema defines the schema for the target connect policy resource.
func (r *targetEnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// DACs do not have an environment
	validTargetTypes := []targettype.TargetType{
		targettype.Bzero,
		targettype.Cluster,
		targettype.Web,
		targettype.Db,
	}

	baseDesc := "Provides management of a BastionZero target's environment. A target can only belong to a single environment."
	resp.Schema = schema.Schema{
		Description: baseDesc,
		MarkdownDescription: baseDesc +
			"\n\nWhen this resource is destroyed, the target's environment is set back to the default environment.",
		Attributes: map[string]schema.Attribute{
			"target_id": schema.StringAttribute{
				Required:    true,
				Description: "The target's ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"target_type": schema.StringAttribute{
				Required:    true,
				Description: fmt.Sprintf("The target's type %s.", internal.PrettyOneOf(validTargetTypes)),
				Validators: []validator.String{
					stringvalidator.OneOf(bastionzero.ToStringSlice(validTargetTypes)...),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: "The environment's ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *targetEnvironmentResource) modifyTargetEnvironment(ctx context.Context, targetType, targetID, environmentID string) diag.Diagnostics {
	var diags diag.Diagnostics

	var err error
	switch targetType {
	case string(targettype.Bzero):
		tflog.Debug(ctx, "Modifying Bzero target's environment")
		_, _, err = r.client.Targets.ModifyBzeroTarget(ctx, targetID, &targets.ModifyBzeroTargetRequest{EnvironmentID: &environmentID})
		break
	case string(targettype.Cluster):
		tflog.Debug(ctx, "Modifying Cluster target's environment")
		_, _, err = r.client.Targets.ModifyClusterTarget(ctx, targetID, &targets.ModifyClusterTargetRequest{EnvironmentID: &environmentID})
		break
	case string(targettype.Web):
		tflog.Debug(ctx, "Modifying Web target's environment")
		_, _, err = r.client.Targets.ModifyWebTarget(ctx, targetID, &targets.ModifyWebTargetRequest{EnvironmentID: &environmentID})
		break
	case string(targettype.Db):
		tflog.Debug(ctx, "Modifying Db target's environment")
		_, _, err = r.client.Targets.ModifyDatabaseTarget(ctx, targetID, &targets.ModifyDatabaseTargetRequest{EnvironmentID: &environmentID})
		break
	default:
		// This should not happen due to validator in schema
		panic(fmt.Sprintf("Unhandled target type: %s. Please report this issue to the provider developers.", targetType))
	}

	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error modifying %v target's environment", targetType),
			fmt.Sprintf("Could not modify %v target's environment, unexpected error: %v", targetType, err.Error()),
		)
		return diags
	}

	return diags
}

func (r *targetEnvironmentResource) modifyTargetEnvironmentWithSchema(ctx context.Context, schema targetEnvironmentModel) diag.Diagnostics {
	targetType := schema.TargetType.ValueString()
	targetID := schema.TargetID.ValueString()
	environmentID := schema.EnvironmentID.ValueString()

	ctx = tflog.SetField(ctx, "target_type", targetType)
	ctx = tflog.SetField(ctx, "target_id", targetID)
	ctx = tflog.SetField(ctx, "environment_id", environmentID)

	return r.modifyTargetEnvironment(ctx, targetType, targetID, environmentID)
}

func (r *targetEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan targetEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.modifyTargetEnvironmentWithSchema(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *targetEnvironmentResource) getEnvironmentID(ctx context.Context, schema targetEnvironmentModel) (string, error) {
	targetType := schema.TargetType.ValueString()
	targetID := schema.TargetID.ValueString()
	environmentID := schema.EnvironmentID.ValueString()
	ctx = tflog.SetField(ctx, "target_type", targetType)
	ctx = tflog.SetField(ctx, "target_id", targetID)
	ctx = tflog.SetField(ctx, "environment_id", environmentID)

	switch targetType {
	case string(targettype.Bzero):
		tflog.Debug(ctx, "Querying for Bzero target's environment")
		bzeroTarget, _, err := r.client.Targets.GetBzeroTarget(ctx, targetID)
		if err != nil {
			return "", err
		}
		return bzeroTarget.EnvironmentID, nil
	case string(targettype.Cluster):
		tflog.Debug(ctx, "Querying for Cluster target's environment")
		clusterTarget, _, err := r.client.Targets.GetClusterTarget(ctx, targetID)
		if err != nil {
			return "", err
		}
		return clusterTarget.EnvironmentID, nil
	case string(targettype.Web):
		tflog.Debug(ctx, "Querying for Web target's environment")
		webTarget, _, err := r.client.Targets.GetWebTarget(ctx, targetID)
		if err != nil {
			return "", err
		}
		return webTarget.EnvironmentID, nil
	case string(targettype.Db):
		tflog.Debug(ctx, "Querying for Db target's environment")
		dbTarget, _, err := r.client.Targets.GetDatabaseTarget(ctx, targetID)
		if err != nil {
			return "", err
		}
		return dbTarget.EnvironmentID, nil
	default:
		// This should not happen due to validator in schema
		panic(fmt.Sprintf("Unhandled target type: %s. Please report this issue to the provider developers.", targetType))
	}
}

// Read refreshes the target connect policy Terraform state with the latest
// data.
func (r *targetEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state targetEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetType := state.TargetType.ValueString()

	refreshedEnvrionmentId, err := r.getEnvironmentID(ctx, state)
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading %v target's environment", targetType),
			fmt.Sprintf("Could not read %v target's environment, unexpected error: %v", targetType, err.Error()),
		)
		return
	}

	state.EnvironmentID = types.StringValue(refreshedEnvrionmentId)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the target connect policy resource and sets the updated Terraform state
// on success.
func (r *targetEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan targetEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.modifyTargetEnvironmentWithSchema(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the target connect policy resource and removes the Terraform state on
// success.
func (r *targetEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state targetEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetType := state.TargetType.ValueString()
	targetID := state.TargetID.ValueString()
	environmentID := state.EnvironmentID.ValueString()
	ctx = tflog.SetField(ctx, "target_type", targetType)
	ctx = tflog.SetField(ctx, "target_id", targetID)
	ctx = tflog.SetField(ctx, "environment_id", environmentID)

	// TODO: Get ID of default environment
	envs, _, err := r.client.Environments.ListEnvironments(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing environments",
			"Could not list environments, unexpected error: "+err.Error(),
		)
		return
	}

	var defaultEnvID string
	for _, env := range envs {
		if env.IsDefault {
			defaultEnvID = env.ID
			break
		}
	}

	// This should never happen due to backend constraint
	if defaultEnvID == "" {
		resp.Diagnostics.AddError(
			"Unexpected error when searching for default environment",
			"Expected to find the default environment, but could not find it. Please report this issue to the provider developers.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found default environment with ID %v", defaultEnvID))
	resp.Diagnostics.Append(r.modifyTargetEnvironment(ctx, targetType, targetID, defaultEnvID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
