package dbtarget

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dbTargetResource{}
	_ resource.ResourceWithConfigure   = &dbTargetResource{}
	_ resource.ResourceWithImportState = &dbTargetResource{}
)

func NewDbTargetResource() resource.Resource {
	return &dbTargetResource{}
}

// dbTargetResource is the resource implementation.
type dbTargetResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *dbTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the db target resource type name.
func (r *dbTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_db_target"
}

// Schema defines the schema for the db target resource.
func (r *dbTargetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a BastionZero database target. Database targets configure remote access to database servers running on [Bzero](bzero_target) targets or [Cluster](cluster_target) targets.",
		Attributes:          makeDbTargetResourceSchema(ctx),
	}
}

// Create creates the db target resource and sets the initial Terraform state.
func (r *dbTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan dbTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createReq := new(targets.CreateDatabaseTargetRequest)
	createReq.TargetName = plan.Name.ValueString()
	createReq.ProxyTargetID = plan.ProxyTargetID.ValueString()
	createReq.RemoteHost = plan.RemoteHost.ValueString()
	// TODO-Yuval: Come back to this. Potential improvement to not include this
	// if its GCP
	createReq.RemotePort = targets.Port{Value: bastionzero.PtrTo(int(plan.RemotePort.ValueInt64()))}
	if !plan.LocalPort.IsNull() {
		createReq.LocalPort = &targets.Port{Value: bastionzero.PtrTo(int(plan.LocalPort.ValueInt64()))}
	}
	createReq.EnvironmentID = plan.EnvironmentID.ValueString()
	createReq.DatabaseAuthenticationConfig = ExpandDatabaseAuthenticationConfig(ctx, plan.DatabaseAuthenticationConfig)

	// Create new db target
	tflog.Debug(ctx, "Creating db target")
	createResp, _, err := r.client.Targets.CreateDatabaseTarget(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating db target",
			"Could not create db target, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "db_target_id", createResp.TargetId)
	tflog.Debug(ctx, "Created db target")
	plan.ID = types.StringValue(createResp.TargetId)

	// Query using the GET API to populate other attributes
	found, diags := readDbTarget(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Failed to find db target after create", "")
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the db target Terraform state with the latest data.
func (r *dbTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state dbTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "db_target_id", state.ID.ValueString())

	// Read db target
	found, diags := readDbTarget(ctx, &state, r.client)
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

// Update updates the db target resource and sets the updated Terraform state on
// success.
func (r *dbTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state dbTargetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "db_target_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modifyReq := new(targets.ModifyDatabaseTargetRequest)
	if !plan.Name.Equal(state.Name) {
		modifyReq.TargetName = bastionzero.PtrTo(plan.Name.ValueString())
	}
	if !plan.ProxyTargetID.Equal(state.ProxyTargetID) {
		modifyReq.ProxyTargetID = bastionzero.PtrTo(plan.ProxyTargetID.ValueString())
	}
	if !plan.RemoteHost.Equal(state.RemoteHost) {
		modifyReq.RemoteHost = bastionzero.PtrTo(plan.RemoteHost.ValueString())
	}
	if !plan.RemotePort.Equal(state.RemotePort) {
		// TODO-Yuval: Come back to this. Potential improvement to not include
		// this if its GCP
		modifyReq.RemotePort = &targets.Port{Value: bastionzero.PtrTo(int(plan.RemotePort.ValueInt64()))}
	}
	if !plan.LocalPort.Equal(state.LocalPort) {
		// TODO-Yuval: The system test needs to check both of these paths
		if !plan.LocalPort.IsNull() {
			modifyReq.LocalPort = &targets.Port{Value: bastionzero.PtrTo(int(plan.LocalPort.ValueInt64()))}
		} else {
			// Send an empty Port value to reset this value to blank
			modifyReq.LocalPort = &targets.Port{Value: nil}
		}
	}
	if !plan.EnvironmentID.Equal(state.EnvironmentID) {
		modifyReq.EnvironmentID = bastionzero.PtrTo(plan.EnvironmentID.ValueString())
	}
	if !plan.DatabaseAuthenticationConfig.Equal(state.DatabaseAuthenticationConfig) {
		modifyReq.DatabaseAuthenticationConfig = ExpandDatabaseAuthenticationConfig(ctx, plan.DatabaseAuthenticationConfig)
	}

	// Update existing db target
	updateResp, _, err := r.client.Targets.ModifyDatabaseTarget(ctx, plan.ID.ValueString(), modifyReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating db target",
			"Could not update db target, unexpected error: "+err.Error(),
		)
		return
	}

	setDbTargetResourceAttributes(ctx, &plan, updateResp)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the db target resource and removes the Terraform state on
// success.
func (r *dbTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dbTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "db_target_id", state.ID.ValueString())

	// Delete existing db target
	tflog.Debug(ctx, "Deleting db target")
	_, err := r.client.Targets.DeleteDatabaseTarget(ctx, state.ID.ValueString())

	// TODO-Yuval: Fix this on backend
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if db target is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting db target",
			"Could not delete db target, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Deleted db target")
}

func (r *dbTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}