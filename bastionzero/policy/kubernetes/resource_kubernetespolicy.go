package kubernetes

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
	_ resource.Resource                     = &kubernetesPolicyResource{}
	_ resource.ResourceWithConfigure        = &kubernetesPolicyResource{}
	_ resource.ResourceWithImportState      = &kubernetesPolicyResource{}
	_ resource.ResourceWithConfigValidators = &kubernetesPolicyResource{}
)

func NewKubernetesPolicyResource() resource.Resource {
	return &kubernetesPolicyResource{}
}

// kubernetesPolicyResource is the resource implementation.
type kubernetesPolicyResource struct {
	client *bastionzero.Client
}

// Configure adds the provider configured BastionZero API client to the
// resource.
func (r *kubernetesPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the Kubernetes policy resource type name.
func (r *kubernetesPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_policy"
}

// Schema defines the schema for the Kubernetes policy resource.
func (r *kubernetesPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a BastionZero Kubernetes policy. Kubernetes policies provide access to Cluster targets.",
		Attributes:  makeKubernetesPolicyResourceSchema(),
	}
}

// Create creates the Kubernetes policy resource and sets the initial Terraform state.
func (r *kubernetesPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan KubernetesPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	p := ExpandKubernetesPolicy(ctx, &plan)

	ctx = tflog.SetField(ctx, "policy_name", p.Name)

	// Create new policy
	tflog.Debug(ctx, "Creating Kubernetes policy")
	createResp, _, err := r.client.Policies.CreateKubernetesPolicy(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Kubernetes policy",
			"Could not create Kubernetes policy, unexpected error: "+err.Error(),
		)
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", createResp.ID)
	tflog.Debug(ctx, "Created Kubernetes policy")

	SetKubernetesPolicyAttributes(ctx, &plan, createResp, false)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Kubernetes policy Terraform state with the latest data.
func (r *kubernetesPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model
	var state KubernetesPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Get refreshed policy value from BastionZero
	tflog.Debug(ctx, "Querying for Kubernetes policy")
	p, _, err := r.client.Policies.GetKubernetesPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// The next terraform plan will recreate the resource
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Kubernetes policy",
			"Could not read Kubernetes policy, unexpected error: "+err.Error())
		return
	}
	tflog.Debug(ctx, "Queried for Kubernetes policy")

	SetKubernetesPolicyAttributes(ctx, &state, p, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the Kubernetes policy resource and sets the updated Terraform state
// on success.
func (r *kubernetesPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan and current state data into the model
	var plan, state KubernetesPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", plan.ID.ValueString())

	// Generate API request body from plan. Only include things in request that
	// have changed between plan and current state
	modPolicy := new(policies.KubernetesPolicy)
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
	if !plan.Clusters.Equal(state.Clusters) {
		modPolicy.Clusters = bastionzero.PtrTo(ExpandPolicyClusters(ctx, plan.Clusters))
	}
	if !plan.ClusterUsers.Equal(state.ClusterUsers) {
		modPolicy.ClusterUsers = bastionzero.PtrTo(ExpandPolicyClusterUsers(ctx, plan.ClusterUsers))
	}
	if !plan.ClusterGroups.Equal(state.ClusterGroups) {
		modPolicy.ClusterGroups = bastionzero.PtrTo(ExpandPolicyClusterGroups(ctx, plan.ClusterGroups))
	}

	// Update existing policy
	updateResp, _, err := r.client.Policies.ModifyKubernetesPolicy(ctx, plan.ID.ValueString(), modPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kubernetes policy",
			"Could not update Kubernetes policy, unexpected error: "+err.Error(),
		)
		return
	}

	SetKubernetesPolicyAttributes(ctx, &plan, updateResp, false)

	// Overwrite with refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the Kubernetes policy resource and removes the Terraform state on
// success.
func (r *kubernetesPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state KubernetesPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "policy_id", state.ID.ValueString())

	// Delete existing policy
	tflog.Debug(ctx, "Deleting Kubernetes policy")
	_, err := r.client.Policies.DeleteKubernetesPolicy(ctx, state.ID.ValueString())
	if apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
		// Return early without error if policy is already deleted
		return
	} else if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Kubernetes policy",
			"Could not delete Kubernetes policy, unexpected error: "+err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted Kubernetes policy")
}

func (r *kubernetesPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *kubernetesPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// Validate that policy is not configured with both environments and
		// clusters (known, non-null values).
		resourcevalidator.Conflicting(
			path.MatchRoot("environments"),
			path.MatchRoot("clusters"),
		),
	}
}
