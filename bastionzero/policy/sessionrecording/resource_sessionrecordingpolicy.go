package sessionrecording

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionRecordingPolicyResource{}
	_ resource.ResourceWithConfigure   = &sessionRecordingPolicyResource{}
	_ resource.ResourceWithImportState = &sessionRecordingPolicyResource{}
)

func NewSessionRecordingPolicyResource() resource.Resource {
	return &sessionRecordingPolicyResource{}
}

// sessionRecordingPolicyResource is the resource implementation.
type sessionRecordingPolicyResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *sessionRecordingPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the session recording policy resource type name.
func (r *sessionRecordingPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sessionrecording_policy"
}

// Schema defines the schema for the session recording policy resource.
func (r *sessionRecordingPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero session recording policy. Session recording policies govern whether users' I/O during shell connections are recorded.",
		Attributes:  makeSessionRecordingPolicyResourceSchema(),
	}
}

// Create creates the session recording policy resource and sets the initial
// Terraform state.
func (r *sessionRecordingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan SessionRecordingPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	p := ExpandSessionRecordingPolicy(ctx, &plan)

	ctx = tflog.SetField(ctx, "policy_name", p.Name)

	// Create new policy
	tflog.Debug(ctx, "Creating session recording policy")
	createResp, _, err := r.client.Policies.CreateSessionRecordingPolicy(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating session recording policy",
			"Could not create session recording policy, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", createResp.ID)
	tflog.Debug(ctx, "Created session recording policy")

	SetSessionRecordingPolicyAttributes(ctx, &plan, createResp, false)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the session recording policy Terraform state with the latest
// data.
func (r *sessionRecordingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state SessionRecordingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Get refreshed policy value from BastionZero
	tflog.Debug(ctx, "Querying for session recording policy")
	p, _, err := r.client.Policies.GetSessionRecordingPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading session recording policy",
			"Could not read session recording policy, unexpected error: "+err.Error())
		return
	}
	tflog.Debug(ctx, "Queried for session recording policy")

	SetSessionRecordingPolicyAttributes(ctx, &state, p, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the session recording policy resource and sets the updated
// Terraform state on success.
func (r *sessionRecordingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state SessionRecordingPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modPolicy := new(policies.SessionRecordingPolicy)
	if !plan.Name.Equal(state.Name) {
		modPolicy.Name = plan.Name.ValueString()
	}
	if !plan.Description.Equal(state.Description) {
		modPolicy.Description = bastionzero.PtrTo(plan.Description.ValueString())
	}
	if !plan.Subjects.Equal(state.Subjects) {
		modPolicy.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, plan.Subjects))
	}
	if !plan.Groups.Equal(state.Groups) {
		modPolicy.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, plan.Groups))
	}
	if !plan.RecordInput.Equal(state.RecordInput) {
		modPolicy.RecordInput = bastionzero.PtrTo(plan.RecordInput.ValueBool())
	}

	// Update existing policy
	updateResp, _, err := r.client.Policies.ModifySessionRecordingPolicy(ctx, plan.ID.ValueString(), modPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating session recording policy",
			"Could not update session recording policy, unexpected error: "+err.Error(),
		)
		return
	}

	SetSessionRecordingPolicyAttributes(ctx, &plan, updateResp, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the session recording policy resource and removes the
// Terraform state on success.
func (r *sessionRecordingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state SessionRecordingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Delete existing policy
	tflog.Debug(ctx, "Deleting session recording policy")
	_, err := r.client.Policies.DeleteSessionRecordingPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if policy is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting session recording policy",
			"Could not delete session recording policy, unexpected error: "+err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted session recording policy")
}

func (r *sessionRecordingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
