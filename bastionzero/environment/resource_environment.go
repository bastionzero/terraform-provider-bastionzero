package environment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

// environmentResource is the resource implementation.
type environmentResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the environment resource type name.
func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the environment resource.
func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero environment resource. An environment is a collection of targets.",
		Attributes:  MakeEnvironmentResourceSchema(),
	}
}

// Create creates the environment resource and sets the initial Terraform state.
func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createReq := new(environments.CreateEnvironmentRequest)
	createReq.Name = plan.Name.ValueString()
	createReq.Description = plan.Description.ValueString()
	createReq.OfflineCleanupTimeoutHours = uint(plan.OfflineCleanupTimeoutHours.ValueInt64())

	ctx = tflog.SetField(ctx, "environment_name", createReq.Name)
	ctx = tflog.SetField(ctx, "environment_desc", createReq.Description)
	ctx = tflog.SetField(ctx, "environment_offline_cleanup_timehout_hours", createReq.OfflineCleanupTimeoutHours)

	// Create new environment
	tflog.Debug(ctx, "Creating environment")
	createResp, _, err := r.client.Environments.CreateEnvironment(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment",
			"Could not create environment, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "environment_id", createResp.ID)
	tflog.Debug(ctx, "Created environment")
	plan.ID = types.StringValue(createResp.ID)

	// Query using the GET API to populate other attributes
	found, diags := readEnvironment(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Failed to find environment after create", "")
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the environment Terraform state with the latest data.
func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

	// Read environment
	found, diags := readEnvironment(ctx, &state, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	}

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the environment resource and sets the updated Terraform state
// on success.
func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state environmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "environment_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modifyReq := new(environments.ModifyEnvironmentRequest)
	if !plan.Description.Equal(state.Description) {
		modifyReq.Description = bastionzero.PtrTo(plan.Description.ValueString())
	}
	if !plan.OfflineCleanupTimeoutHours.Equal(state.OfflineCleanupTimeoutHours) {
		modifyReq.OfflineCleanupTimeoutHours = bastionzero.PtrTo(uint(plan.OfflineCleanupTimeoutHours.ValueInt64()))
	}

	// Update existing environment
	_, err := r.client.Environments.ModifyEnvironment(ctx, plan.ID.ValueString(), modifyReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment",
			"Could not update environment, unexpected error: "+err.Error(),
		)
		return
	}

	// Query using the GET API to populate other attributes
	found, diags := readEnvironment(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Failed to find environment after update", "")
		return
	}

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the environment resource and removes the Terraform state on
// success.
func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

	// Present user-friendly error instead of internal server error if
	// environment contains targets
	if len(state.Targets.Elements()) > 0 {
		resp.Diagnostics.AddError(
			"Error deleting environment",
			fmt.Sprintf("Cannot delete an environment with targets in it. Environment %s contains %d target(s).\n", state.ID.ValueString(), len(state.Targets.Elements()))+
				"Please remove all targets from this environment before destroying.",
		)
		return
	}

	// Delete existing environment
	tflog.Debug(ctx, "Deleting environment")
	_, err := r.client.Environments.DeleteEnvironment(ctx, state.ID.ValueString())

	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if environment is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting environment",
			"Could not delete environment, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Deleted environment")
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
