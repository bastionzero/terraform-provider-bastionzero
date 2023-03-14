package targetconnect

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/typesext"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &targetConnectPolicyResource{}
	_ resource.ResourceWithConfigure   = &targetConnectPolicyResource{}
	_ resource.ResourceWithImportState = &targetConnectPolicyResource{}
)

func NewTargetConnectPolicyResource() resource.Resource {
	return &targetConnectPolicyResource{}
}

// targetConnectPolicyModel maps the target connect policy schema data.
type targetConnectPolicyModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Subjects     types.Set    `tfsdk:"subjects"`
	Groups       types.Set    `tfsdk:"groups"`
	Environments types.Set    `tfsdk:"environments"`
	Targets      types.Set    `tfsdk:"targets"`
	TargetUsers  types.Set    `tfsdk:"target_users"`
	Verbs        types.Set    `tfsdk:"verbs"`
}

// targetConnectPolicyResource is the resource implementation.
type targetConnectPolicyResource struct {
	client *bastionzero.Client
}

// setTargetConnectPolicyAttributes populates the TF schema data from a target
// connect policy
func setTargetConnectPolicyAttributes(ctx context.Context, schema *targetConnectPolicyModel, apiPolicy *policies.TargetConnectPolicy) {
	schema.ID = types.StringValue(apiPolicy.ID)
	schema.Name = types.StringValue(apiPolicy.Name)
	schema.Description = typesext.StringEmptyIsNullValue(apiPolicy.Description)
	schema.Subjects = policy.FlattenPolicySubjects(ctx, apiPolicy.Subjects)
	schema.Groups = policy.FlattenPolicyGroups(ctx, apiPolicy.Groups)
	schema.Environments = policy.FlattenPolicyEnvironments(ctx, apiPolicy.Environments)
	schema.Targets = policy.FlattenPolicyTargets(ctx, apiPolicy.Targets)
	schema.TargetUsers = FlattenPolicyTargetUsers(ctx, apiPolicy.TargetUsers)
	schema.Verbs = FlattenPolicyVerbs(ctx, apiPolicy.Verbs)
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
func (r *targetConnectPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

	// Prevent empty string for policy name
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty policy name",
			"All policies require a name. Set the name value in the configuration. Do not use an empty string.",
		)
		return
	}

	// Generate API request body from plan
	p := new(policies.TargetConnectPolicy)
	p.Name = plan.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, plan.Description)
	p.Subjects = policy.ExpandPolicySubjects(ctx, plan.Subjects)
	p.Groups = policy.ExpandPolicyGroups(ctx, plan.Groups)
	p.Environments = policy.ExpandPolicyEnvironments(ctx, plan.Environments)
	p.Targets = policy.ExpandPolicyTargets(ctx, plan.Targets)
	p.TargetUsers = ExpandPolicyTargetUsers(ctx, plan.TargetUsers)
	p.Verbs = ExpandPolicyVerbs(ctx, plan.Verbs)

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

	setTargetConnectPolicyAttributes(ctx, &plan, createResp)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the target connect policy Terraform state with the latest data.
func (r *targetConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// // Read Terraform prior state data into the model
	// var state environmentModel
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

	// // Read environment
	// found, diags := readEnvironment(ctx, &state, r.client)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// if !found {
	// 	// The next terraform plan will recreate thre resource
	// 	resp.State.RemoveResource(ctx)
	// 	return
	// }

	// // Overwrite with refreshed state
	// resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the target connect policy resource and sets the updated Terraform state
// on success.
func (r *targetConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// // Read Terraform plan and current state data into the model
	// var plan, state environmentModel

	// resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// ctx = tflog.SetField(ctx, "environment_id", plan.ID.ValueString())

	// // Generate API request body from plan. Only include things in request that
	// // have changed between plan and current state
	// modifyReq := new(environments.ModifyEnvironmentRequest)
	// if !plan.Description.Equal(state.Description) {
	// 	modifyReq.Description = bastionzero.PtrTo(plan.Description.ValueString())
	// }
	// if !plan.OfflineCleanupTimeoutHours.Equal(state.OfflineCleanupTimeoutHours) {
	// 	modifyReq.OfflineCleanupTimeoutHours = bastionzero.PtrTo(uint(plan.OfflineCleanupTimeoutHours.ValueInt64()))
	// }

	// // Update existing environment
	// _, err := r.client.Environments.ModifyEnvironment(ctx, plan.ID.ValueString(), modifyReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating environment",
	// 		"Could not update environment, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// // Query using the GET API to populate other attributes
	// found, diags := readEnvironment(ctx, &plan, r.client)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// if !found {
	// 	resp.Diagnostics.AddError("Failed to find environment after update", "")
	// 	return
	// }

	// // Overwrite with refreshed state
	// resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the target connect policy resource and removes the Terraform state on
// success.
func (r *targetConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// // Retrieve values from state
	// var state environmentModel
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

	// // Present user-friendly error instead of internal server error if
	// // environment contains targets
	// if len(state.Targets.Elements()) > 0 {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting environment",
	// 		fmt.Sprintf("Cannot delete an environment with targets in it. Environment %s contains %d target(s).\n", state.ID.ValueString(), len(state.Targets.Elements()))+
	// 			"Please remove all targets from this environment before destroying.",
	// 	)
	// 	return
	// }

	// // Delete existing environment
	// tflog.Debug(ctx, "Deleting environment")
	// httpResp, err := r.client.Environments.DeleteEnvironment(ctx, state.ID.ValueString())
	// if httpResp.StatusCode == http.StatusNotFound {
	// 	// Return early without error if environment is already deleted
	// 	return
	// }
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting environment",
	// 		"Could not delete environment, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
	// tflog.Debug(ctx, "Deleted environment")
}

func (r *targetConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
