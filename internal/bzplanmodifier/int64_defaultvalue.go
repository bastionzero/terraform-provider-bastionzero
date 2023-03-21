package bzplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64DefaultValue is a plan modifier that injects a default Int64 value into
// the plan if the configuration contains a null value.
func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &int64DefaultValuePlanModifier{v}
}

type int64DefaultValuePlanModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*int64DefaultValuePlanModifier)(nil)

// Description returns a plain text description of the validator's behavior,
// suitable for a practitioner to understand its impact.
func (m *int64DefaultValuePlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.DefaultValue)
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (m *int64DefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.DefaultValue)
}

// PlanModifyInt64 runs the logic of the plan modifier. Access to the
// configuration, plan, and state is available in `req`, while `resp` contains
// fields for updating the planned value, triggering resource replacement, and
// returning diagnostics.
func (m *int64DefaultValuePlanModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	res.PlanValue = m.DefaultValue
}
