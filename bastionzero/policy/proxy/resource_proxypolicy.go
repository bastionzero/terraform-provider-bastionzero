package proxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/apierror"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &proxyPolicyResource{}
	_ resource.ResourceWithConfigure        = &proxyPolicyResource{}
	_ resource.ResourceWithImportState      = &proxyPolicyResource{}
	_ resource.ResourceWithConfigValidators = &proxyPolicyResource{}
)

func NewProxyPolicyResource() resource.Resource {
	return &proxyPolicyResource{}
}

// proxyPolicyResource is the resource implementation.
type proxyPolicyResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *proxyPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the proxy policy resource type name.
func (r *proxyPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_policy"
}

// Schema defines the schema for the proxy policy resource.
func (r *proxyPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero proxy policy. Proxy policies provide access to Db and Web targets.",
		Attributes:  makeProxyPolicyResourceSchema(),
	}
}

// Create creates the proxy policy resource and sets the initial Terraform
// state.
func (r *proxyPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan ProxyPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	p := ExpandProxyPolicy(ctx, &plan)

	ctx = tflog.SetField(ctx, "policy_name", p.Name)

	// Create new policy
	tflog.Debug(ctx, "Creating proxy policy")
	createResp, _, err := r.client.Policies.CreateProxyPolicy(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating proxy policy",
			"Could not create proxy policy, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", createResp.ID)
	tflog.Debug(ctx, "Created proxy policy")

	SetProxyPolicyAttributes(ctx, &plan, createResp, false)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the proxy policy Terraform state with the latest data.
func (r *proxyPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state ProxyPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Get refreshed policy value from BastionZero
	tflog.Debug(ctx, "Querying for proxy policy")
	p, _, err := r.client.Policies.GetProxyPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading proxy policy",
			"Could not read proxy policy, unexpected error: "+err.Error())
		return
	}
	tflog.Debug(ctx, "Queried for proxy policy")

	SetProxyPolicyAttributes(ctx, &state, p, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the proxy policy resource and sets the updated Terraform state
// on success.
func (r *proxyPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state ProxyPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modPolicy := new(policies.ProxyPolicy)
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
		modPolicy.TargetUsers = bastionzero.PtrTo(policy.ExpandPolicyTargetUsers(ctx, plan.TargetUsers))
	}

	// Update existing policy
	updateResp, _, err := r.client.Policies.ModifyProxyPolicy(ctx, plan.ID.ValueString(), modPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating proxy policy",
			"Could not update proxy policy, unexpected error: "+err.Error(),
		)
		return
	}

	SetProxyPolicyAttributes(ctx, &plan, updateResp, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the proxy policy resource and removes the Terraform state on
// success.
func (r *proxyPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ProxyPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Delete existing policy
	tflog.Debug(ctx, "Deleting proxy policy")
	_, err := r.client.Policies.DeleteProxyPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if policy is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting proxy policy",
			"Could not delete proxy policy, unexpected error: "+err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted proxy policy")
}

func (r *proxyPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *proxyPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// Validate that policy is not configured with both environments and
		// targets (known, non-null values).
		resourcevalidator.Conflicting(
			path.MatchRoot("environments"),
			path.MatchRoot("targets"),
		),
	}
}
