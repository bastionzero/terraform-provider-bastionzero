package typesext

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Int64PointerValue creates an Int64 with a null value if nil or a known value.
func Int64PointerValue(value *int) basetypes.Int64Value {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*value))
}
