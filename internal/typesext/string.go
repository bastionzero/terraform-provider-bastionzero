package typesext

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// StringEmptyIsNullValue creates a String with a null value if value is nil or
// if value is equal to the empty string. Otherwise, a known value.
func StringEmptyIsNullValue(value *string) basetypes.StringValue {
	if value == nil || *value == "" {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(*value)
}
