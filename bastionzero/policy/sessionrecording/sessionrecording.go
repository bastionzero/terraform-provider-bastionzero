package sessionrecording

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SessionRecordingPolicyModel maps the session recording policy schema data.
type SessionRecordingPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Subjects    types.Set    `tfsdk:"subjects"`
	Groups      types.Set    `tfsdk:"groups"`
	RecordInput types.Bool   `tfsdk:"record_input"`
}

func (m *SessionRecordingPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *SessionRecordingPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *SessionRecordingPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *SessionRecordingPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *SessionRecordingPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *SessionRecordingPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *SessionRecordingPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *SessionRecordingPolicyModel) GetGroups() types.Set   { return m.Groups }

// SetSessionRecordingPolicyAttributes populates the TF schema data from a
// session recording policy
func SetSessionRecordingPolicyAttributes(ctx context.Context, schema *SessionRecordingPolicyModel, apiPolicy *policies.SessionRecordingPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)
	schema.RecordInput = types.BoolValue(apiPolicy.GetRecordInput())
}

func ExpandSessionRecordingPolicy(ctx context.Context, schema *SessionRecordingPolicyModel) *policies.SessionRecordingPolicy {
	p := new(policies.SessionRecordingPolicy)
	p.Name = schema.Name.ValueString()
	p.Description = internal.StringFromFramework(ctx, schema.Description)
	p.Subjects = bastionzero.PtrTo(policy.ExpandPolicySubjects(ctx, schema.Subjects))
	p.Groups = bastionzero.PtrTo(policy.ExpandPolicyGroups(ctx, schema.Groups))
	p.RecordInput = bastionzero.PtrTo(schema.RecordInput.ValueBool())

	return p
}

func makeSessionRecordingPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.SessionRecording)
	attributes["record_input"] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "If `true`, then in addition to session output, session input should be recorded. If `false`, then only session output should be recorded (Defaults to `false`).",
		// Don't allow null value to make it easier when parsing results back
		// into TF
		Default: booldefault.StaticBool(false),
	}

	return attributes
}
