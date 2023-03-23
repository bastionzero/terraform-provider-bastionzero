package bzplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BoolDefaultValue is a plan modifier that injects a default Bool value into
// the plan if the configuration contains a null value.
func BoolDefaultValue(v types.Bool) planmodifier.Bool {
	return &boolDefaultValueProvider{v}
}

type boolDefaultValueProvider struct {
	DefaultValue types.Bool
}

var _ planmodifier.Bool = (*boolDefaultValueProvider)(nil)

// Description returns a plain text description of the validator's behavior,
// suitable for a practitioner to understand its impact.
func (m *boolDefaultValueProvider) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.DefaultValue)
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (m *boolDefaultValueProvider) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.DefaultValue)
}

// PlanModifyInt64 runs the logic of the plan modifier. Access to the
// configuration, plan, and state is available in `req`, while `resp` contains
// fields for updating the planned value, triggering resource replacement, and
// returning diagnostics.
func (m *boolDefaultValueProvider) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	res.PlanValue = m.DefaultValue
}