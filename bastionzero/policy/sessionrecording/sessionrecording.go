package sessionrecording

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sessionRecordingPolicyModel maps the session recording policy schema data.
type sessionRecordingPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Subjects    types.Set    `tfsdk:"subjects"`
	Groups      types.Set    `tfsdk:"groups"`
	RecordInput types.Bool   `tfsdk:"record_input"`
}

func (m *sessionRecordingPolicyModel) SetID(value types.String)          { m.ID = value }
func (m *sessionRecordingPolicyModel) SetName(value types.String)        { m.Name = value }
func (m *sessionRecordingPolicyModel) SetType(value types.String)        { m.Type = value }
func (m *sessionRecordingPolicyModel) SetDescription(value types.String) { m.Description = value }
func (m *sessionRecordingPolicyModel) SetSubjects(value types.Set)       { m.Subjects = value }
func (m *sessionRecordingPolicyModel) SetGroups(value types.Set)         { m.Groups = value }

func (m *sessionRecordingPolicyModel) GetSubjects() types.Set { return m.Subjects }
func (m *sessionRecordingPolicyModel) GetGroups() types.Set   { return m.Groups }

// setSessionRecordingPolicyAttributes populates the TF schema data from a
// session recording policy
func setSessionRecordingPolicyAttributes(ctx context.Context, schema *sessionRecordingPolicyModel, apiPolicy *policies.SessionRecordingPolicy, modelIsDataSource bool) {
	policy.SetBasePolicyAttributes(ctx, schema, apiPolicy, modelIsDataSource)
	schema.RecordInput = types.BoolValue(apiPolicy.GetRecordInput())
}

func makeSessionRecordingPolicyResourceSchema() map[string]schema.Attribute {
	attributes := policy.BasePolicyResourceAttributes(policytype.SessionRecording)
	attributes["record_input"] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "If true, then in addition to session output, session input should be recorded. If false, then only session output should be recorded (Defaults to false).",
		// Don't allow null value to make it easier when parsing results back
		// into TF
		Default: booldefault.StaticBool(false),
	}

	return attributes
}
