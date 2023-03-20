package typesext

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NumberPointerValue creates a Number with a null value if nil or a known
// value.
func NumberPointerValue(value *int) basetypes.NumberValue {
	if value == nil {
		return types.NumberNull()
	}

	return types.NumberValue(new(big.Float).SetInt(big.NewInt(int64(*value))))
}
