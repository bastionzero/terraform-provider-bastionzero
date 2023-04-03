package jit

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	_ resource.Resource                = &jitPolicyResource{}
	_ resource.ResourceWithConfigure   = &jitPolicyResource{}
	_ resource.ResourceWithImportState = &jitPolicyResource{}
)

func NewJITPolicyResource() resource.Resource {
	return &jitPolicyResource{}
}

// jitPolicyResource is the resource implementation.
type jitPolicyResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *jitPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the JIT policy resource type name.
func (r *jitPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jit_policy"
}

// Schema defines the schema for the JIT policy resource.
func (r *jitPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero JIT policy. JIT policies provide just in time access to targets." +
			fmt.Sprintf("\n\n~> **Note on child policies** A JIT policy's `child_policies` can only refer to policies of the following types: %v. If any of the referenced policies ", strings.Join(bastionzero.ToStringSlice(allowedChildPolicyTypes()), ", ")) +
			"are not of the valid type, then an error is returned when creating/updating the JIT policy.",
		Attributes: makeJITPolicyResourceSchema(),
	}
}

// Create creates the JIT policy resource and sets the initial Terraform state.
func (r *jitPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan JITPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	p := ExpandJITPolicy(ctx, &plan)

	ctx = tflog.SetField(ctx, "policy_name", p.Name)

	// Create new policy
	tflog.Debug(ctx, "Creating JIT policy")
	createResp, _, err := r.client.Policies.CreateJITPolicy(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating JIT policy",
			"Could not create JIT policy, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", createResp.ID)
	tflog.Debug(ctx, "Created JIT policy")

	SetJITPolicyAttributes(ctx, &plan, createResp, false)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the JIT policy Terraform state with the latest data.
func (r *jitPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state JITPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Get refreshed policy value from BastionZero
	tflog.Debug(ctx, "Querying for JIT policy")
	p, _, err := r.client.Policies.GetJITPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading JIT policy",
			"Could not read JIT policy, unexpected error: "+err.Error())
		return
	}
	tflog.Debug(ctx, "Queried for JIT policy")

	SetJITPolicyAttributes(ctx, &state, p, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the JIT policy resource and sets the updated Terraform state
// on success.
func (r *jitPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state JITPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modPolicy := new(policies.ModifyJITPolicyRequest)
	if !plan.Name.Equal(state.Name) {
		modPolicy.Name = bastionzero.PtrTo(plan.Name.ValueString())
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

	// Must always provide child policies
	modPolicy.ChildPolicies = bastionzero.PtrTo(ExpandChildPolicies(ctx, plan.ChildPolicies))

	if !plan.AutomaticallyApproved.Equal(state.AutomaticallyApproved) {
		modPolicy.AutomaticallyApproved = bastionzero.PtrTo(plan.AutomaticallyApproved.ValueBool())
	}
	if !plan.Duration.Equal(state.Duration) {
		modPolicy.Duration = bastionzero.PtrTo(uint(plan.Duration.ValueInt64()))
	}

	// Update existing policy
	updateResp, _, err := r.client.Policies.ModifyJITPolicy(ctx, plan.ID.ValueString(), modPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating JIT policy",
			"Could not update JIT policy, unexpected error: "+err.Error(),
		)
		return
	}

	SetJITPolicyAttributes(ctx, &plan, updateResp, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the JIT policy resource and removes the Terraform state on
// success.
func (r *jitPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state JITPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Delete existing policy
	tflog.Debug(ctx, "Deleting JIT policy")
	_, err := r.client.Policies.DeleteJITPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if policy is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting JIT policy",
			"Could not delete JIT policy, unexpected error: "+err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted JIT policy")
}

func (r *jitPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
