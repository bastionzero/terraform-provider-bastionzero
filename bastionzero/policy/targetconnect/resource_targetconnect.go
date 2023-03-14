package targetconnect

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
// 	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
// 	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
// 	"github.com/hashicorp/terraform-plugin-framework/attr"
// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// )

// // Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource                = &targetConnectPolicyResource{}
// 	_ resource.ResourceWithConfigure   = &targetConnectPolicyResource{}
// 	_ resource.ResourceWithImportState = &targetConnectPolicyResource{}
// )

// func NewTargetConnectPolicyResource() resource.Resource {
// 	return &targetConnectPolicyResource{}
// }

// // targetConnectPolicyModel maps the target connect policy schema data.
// type targetConnectPolicyModel struct {
// 	ID          types.String `tfsdk:"id"`
// 	Name        types.String `tfsdk:"name"`
// 	Description types.String `tfsdk:"description"`

// 	// Subjects TODO

// 	TimeCreated                types.String `tfsdk:"time_created"`
// 	OfflineCleanupTimeoutHours types.Int64  `tfsdk:"offline_cleanup_timeout_hours"`
// 	Targets                    types.Map    `tfsdk:"targets"` // key is target id. value is environmentTargetModel
// }

// // environmentTargetModel maps target summary data.
// type environmentTargetModel struct {
// 	ID   types.String `tfsdk:"id"`
// 	Type types.String `tfsdk:"type"`
// }

// var (
// 	environmentTargetModelAttrTypes = map[string]attr.Type{
// 		"id":   types.StringType,
// 		"type": types.StringType,
// 	}
// )

// // targetConnectPolicyResource is the resource implementation.
// type targetConnectPolicyResource struct {
// 	client *bastionzero.Client
// }

// // setTargetConnectPolicyAttributes populates the TF schema data from a target connect policy
// func setTargetConnectPolicyAttributes(ctx context.Context, schema *targetConnectPolicyModel, policy *policies.TargetConnectPolicy) diag.Diagnostics {
// 	schema.Name = types.StringValue(env.Name)

// 	// Use StringEmptyIsNullValue to normalize "" to Terraform Null value (since
// 	// the schema says description is optional/nullable)
// 	// schema.Description = typesext.StringEmptyIsNullValue(&env.Description)

// 	// Preserve null in TF schema. We say that "" is semantically equivalent to
// 	// null for the environment schema
// 	if schema.Description.IsNull() && env.Description == "" {
// 		schema.Description = types.StringNull()
// 	} else {
// 		schema.Description = types.StringValue(env.Description)
// 	}

// 	schema.OfflineCleanupTimeoutHours = types.Int64Value(int64(env.OfflineCleanupTimeoutHours))

// 	schema.ID = types.StringValue(env.ID)
// 	schema.OrganizationID = types.StringValue(env.OrganizationID)
// 	schema.IsDefault = types.BoolValue(env.IsDefault)
// 	schema.TimeCreated = types.StringValue(env.TimeCreated.UTC().Format(time.RFC3339))

// 	targetsMap := make(map[string]environmentTargetModel)
// 	for _, target := range env.Targets {
// 		targetsMap[target.ID] = environmentTargetModel{
// 			ID:   types.StringValue(target.ID),
// 			Type: types.StringValue(string(target.Type)),
// 		}
// 	}

// 	types.ListValue()
// 	targets, diags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: environmentTargetModelAttrTypes}, targetsMap)
// 	schema.Targets = targets

// 	types.SetValueMust()

// 	return diags
// }

// // Configure adds the provider configured BastionZero API client to the
// // resource.
// func (r *targetConnectPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	// Prevent panic if the provider has not been configured.
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(*bastionzero.Client)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource configure type",
// 			fmt.Sprintf("Expected *bastionzero.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.client = client
// }

// // Metadata returns the target connect policy resource type name.
// func (r *targetConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_targetconnect_policy"
// }

// // Schema defines the schema for the target connect policy resource.
// func (r *targetConnectPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Description: "Provides a BastionZero target connect policy. Target connect policies provide access to Bzero and DynamicAccessConfig targets.",
// 		// Attributes:  makeEnvironmentResourceSchema(),
// 	}
// }

// // Create creates the target connect policy resource and sets the initial Terraform state.
// func (r *targetConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	// Read Terraform plan data into the model
// 	var plan environmentModel
// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Prevent empty string for environment name
// 	if plan.Name.ValueString() == "" {
// 		resp.Diagnostics.AddAttributeError(
// 			path.Root("name"),
// 			"Empty environment name",
// 			"All environments require a name. Set the name value in the configuration. Do not use an empty string.",
// 		)
// 		return
// 	}

// 	// Generate API request body from plan
// 	createReq := new(environments.CreateEnvironmentRequest)
// 	createReq.Name = plan.Name.ValueString()
// 	createReq.Description = plan.Description.ValueString()
// 	createReq.OfflineCleanupTimeoutHours = uint(plan.OfflineCleanupTimeoutHours.ValueInt64())

// 	ctx = tflog.SetField(ctx, "environment_name", createReq.Name)
// 	ctx = tflog.SetField(ctx, "environment_desc", createReq.Description)
// 	ctx = tflog.SetField(ctx, "environment_offline_cleanup_timehout_hours", createReq.OfflineCleanupTimeoutHours)

// 	// Create new environment
// 	tflog.Debug(ctx, "Creating environment")
// 	createResp, _, err := r.client.Environments.CreateEnvironment(ctx, createReq)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error creating environment",
// 			"Could not create environment, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}
// 	ctx = tflog.SetField(ctx, "environment_id", createResp.ID)
// 	tflog.Debug(ctx, "Created environment")
// 	plan.ID = types.StringValue(createResp.ID)

// 	// Query using the GET API to populate other attributes
// 	found, diags := readEnvironment(ctx, &plan, r.client)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	if !found {
// 		resp.Diagnostics.AddError("Failed to find environment after create", "")
// 		return
// 	}

// 	// Set state to fully populated data
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
// }

// // Read refreshes the target connect policy Terraform state with the latest data.
// func (r *targetConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	// Read Terraform prior state data into the model
// 	var state environmentModel
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

// 	// Read environment
// 	found, diags := readEnvironment(ctx, &state, r.client)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	if !found {
// 		// The next terraform plan will recreate thre resource
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}

// 	// Overwrite with refreshed state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
// }

// // Update updates the target connect policy resource and sets the updated Terraform state
// // on success.
// func (r *targetConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	// Read Terraform plan and current state data into the model
// 	var plan, state environmentModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	ctx = tflog.SetField(ctx, "environment_id", plan.ID.ValueString())

// 	// Generate API request body from plan. Only include things in request that
// 	// have changed between plan and current state
// 	modifyReq := new(environments.ModifyEnvironmentRequest)
// 	if !plan.Description.Equal(state.Description) {
// 		modifyReq.Description = bastionzero.PtrTo(plan.Description.ValueString())
// 	}
// 	if !plan.OfflineCleanupTimeoutHours.Equal(state.OfflineCleanupTimeoutHours) {
// 		modifyReq.OfflineCleanupTimeoutHours = bastionzero.PtrTo(uint(plan.OfflineCleanupTimeoutHours.ValueInt64()))
// 	}

// 	// Update existing environment
// 	_, err := r.client.Environments.ModifyEnvironment(ctx, plan.ID.ValueString(), modifyReq)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error updating environment",
// 			"Could not update environment, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}

// 	// Query using the GET API to populate other attributes
// 	found, diags := readEnvironment(ctx, &plan, r.client)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	if !found {
// 		resp.Diagnostics.AddError("Failed to find environment after update", "")
// 		return
// 	}

// 	// Overwrite with refreshed state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
// }

// // Delete deletes the target connect policy resource and removes the Terraform state on
// // success.
// func (r *targetConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	// Retrieve values from state
// 	var state environmentModel
// 	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	ctx = tflog.SetField(ctx, "environment_id", state.ID.ValueString())

// 	// Present user-friendly error instead of internal server error if
// 	// environment contains targets
// 	if len(state.Targets.Elements()) > 0 {
// 		resp.Diagnostics.AddError(
// 			"Error deleting environment",
// 			fmt.Sprintf("Cannot delete an environment with targets in it. Environment %s contains %d target(s).\n", state.ID.ValueString(), len(state.Targets.Elements()))+
// 				"Please remove all targets from this environment before destroying.",
// 		)
// 		return
// 	}

// 	// Delete existing environment
// 	tflog.Debug(ctx, "Deleting environment")
// 	httpResp, err := r.client.Environments.DeleteEnvironment(ctx, state.ID.ValueString())
// 	if httpResp.StatusCode == http.StatusNotFound {
// 		// Return early without error if environment is already deleted
// 		return
// 	}
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error deleting environment",
// 			"Could not delete environment, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}
// 	tflog.Debug(ctx, "Deleted environment")
// }

// func (r *targetConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
