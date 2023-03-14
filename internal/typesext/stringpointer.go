package typesext

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// StringPointerValue creates a String with a null value if nil or a known
// value.
func StringPointerValue(value *string) basetypes.StringValue {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(*value)
}
