package targetconnect

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &targetConnectPolicyResource{}
	_ resource.ResourceWithConfigure        = &targetConnectPolicyResource{}
	_ resource.ResourceWithImportState      = &targetConnectPolicyResource{}
	_ resource.ResourceWithConfigValidators = &targetConnectPolicyResource{}
)

func NewTargetConnectPolicyResource() resource.Resource {
	return &targetConnectPolicyResource{}
}

// targetConnectPolicyResource is the resource implementation.
type targetConnectPolicyResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *targetConnectPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *targetConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_targetconnect_policy"
}

// Schema defines the schema for the target connect policy resource.
func (r *targetConnectPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero target connect policy. Target connect policies provide access to Bzero and DynamicAccessConfig targets.",
		Attributes:  makeTargetConnectPolicyResourceSchema(),
	}
}

// Create creates the target connect policy resource and sets the initial Terraform state.
func (r *targetConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan targetConnectPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	p := new(policies.TargetConnectPolicy)
	p.Name = plan.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, plan.Description)
	p.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, plan.Subjects))
	p.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, plan.Groups))
	p.Environments = bastionzero.PtrTo(policy.ExpandPolicyEnvironments(ctx, plan.Environments))
	p.Targets = bastionzero.PtrTo(policy.ExpandPolicyTargets(ctx, plan.Targets))
	p.TargetUsers = bastionzero.PtrTo(ExpandPolicyTargetUsers(ctx, plan.TargetUsers))
	p.Verbs = bastionzero.PtrTo(ExpandPolicyVerbs(ctx, plan.Verbs))

	ctx = tflog.SetField(ctx, "policy_name", p.Name)

	// Create new policy
	tflog.Debug(ctx, "Creating target connect policy")
	createResp, _, err := r.client.Policies.CreateTargetConnectPolicy(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating target connect policy",
			"Could not create target connect policy, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", createResp.ID)
	tflog.Debug(ctx, "Created target connect policy")

	setTargetConnectPolicyAttributes(ctx, &plan, createResp, false)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the target connect policy Terraform state with the latest data.
func (r *targetConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state targetConnectPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Get refreshed environment value from BastionZero
	tflog.Debug(ctx, "Querying for target connect policy")
	p, _, err := r.client.Policies.GetTargetConnectPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading target connect policy",
			"Could not read target connect policy, unexpected error: "+err.Error())
		return
	}
	tflog.Debug(ctx, "Queried for target connect policy")

	setTargetConnectPolicyAttributes(ctx, &state, p, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the target connect policy resource and sets the updated Terraform state
// on success.
func (r *targetConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state targetConnectPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modPolicy := new(policies.TargetConnectPolicy)
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
	if !plan.Environments.Equal(state.Environments) {
		modPolicy.Environments = bastionzero.PtrTo(policy.ExpandPolicyEnvironments(ctx, plan.Environments))
	}
	if !plan.Targets.Equal(state.Targets) {
		modPolicy.Targets = bastionzero.PtrTo(policy.ExpandPolicyTargets(ctx, plan.Targets))
	}
	if !plan.TargetUsers.Equal(state.TargetUsers) {
		modPolicy.TargetUsers = bastionzero.PtrTo(ExpandPolicyTargetUsers(ctx, plan.TargetUsers))
	}
	if !plan.Verbs.Equal(state.Verbs) {
		modPolicy.Verbs = bastionzero.PtrTo(ExpandPolicyVerbs(ctx, plan.Verbs))
	}

	// Update existing policy
	updateResp, _, err := r.client.Policies.ModifyTargetConnectPolicy(ctx, plan.ID.ValueString(), modPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating target connect policy",
			"Could not update target connect policy, unexpected error: "+err.Error(),
		)
		return
	}

	setTargetConnectPolicyAttributes(ctx, &plan, updateResp, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the target connect policy resource and removes the Terraform state on
// success.
func (r *targetConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state targetConnectPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Delete existing environment
	tflog.Debug(ctx, "Deleting target connect policy")
	_, err := r.client.Policies.DeleteTargetConnectPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if policy is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting target connect policy",
			"Could not delete target connect policy, unexpected error: "+err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted target connect policy")
}

func (r *targetConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *targetConnectPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// Validate that policy is not configured with both environments and
		// targets (known, non-null values).
		resourcevalidator.Conflicting(
			path.MatchRoot("environments"),
			path.MatchRoot("targets"),
		),
	}
}
