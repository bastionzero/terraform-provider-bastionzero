package planmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringDefaultValue is a plan modifier that injects a default String value
// into the plan if the configuration contains a null value.
func StringDefaultValue(v types.String) planmodifier.String {
	return &stringDefaultValueProvider{v}
}

type stringDefaultValueProvider struct {
	DefaultValue types.String
}

var _ planmodifier.String = (*stringDefaultValueProvider)(nil)

// Description returns a plain text description of the validator's behavior,
// suitable for a practitioner to understand its impact.
func (m *stringDefaultValueProvider) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.DefaultValue)
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (m *stringDefaultValueProvider) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.DefaultValue)
}

// PlanModifyInt64 runs the logic of the plan modifier. Access to the
// configuration, plan, and state is available in `req`, while `resp` contains
// fields for updating the planned value, triggering resource replacement, and
// returning diagnostics.
func (m *stringDefaultValueProvider) PlanModifyString(_ context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	res.PlanValue = m.DefaultValue
}
