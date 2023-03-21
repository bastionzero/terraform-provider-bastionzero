package bzplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetDefaultValue is a plan modifier that injects a default Set value into the
// plan if the configuration contains a null value.
func SetDefaultValue(v types.Set) planmodifier.Set {
	return &setDefaultValueProvider{v}
}

type setDefaultValueProvider struct {
	DefaultValue types.Set
}

var _ planmodifier.Set = (*setDefaultValueProvider)(nil)

// Description returns a plain text description of the validator's behavior,
// suitable for a practitioner to understand its impact.
func (m *setDefaultValueProvider) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.DefaultValue)
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (m *setDefaultValueProvider) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.DefaultValue)
}

// PlanModifyInt64 runs the logic of the plan modifier. Access to the
// configuration, plan, and state is available in `req`, while `resp` contains
// fields for updating the planned value, triggering resource replacement, and
// returning diagnostics.
func (m *setDefaultValueProvider) PlanModifySet(_ context.Context, req planmodifier.SetRequest, res *planmodifier.SetResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	res.PlanValue = m.DefaultValue
}
