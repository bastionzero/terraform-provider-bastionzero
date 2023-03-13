package typesext

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type String = StringValue

// NewStringPointerValue creates a String with a null value if nil or a known
// value.
func NewStringPointerValue(value *string) StringValue {
	if value == nil {
		return StringValue{types.StringNull()}
	}

	return StringValue{types.StringValue(*value)}
}

// StringPointerValue creates a String with a null value if nil or a known
// value.
func StringPointerValue(value *string) StringValue {
	return NewStringPointerValue(value)
}

type StringType struct {
	basetypes.StringType
}

func (c StringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := c.StringType.ValueFromTerraform(ctx, in)

	return StringValue{
		// unchecked type assertion
		val.(basetypes.StringValue),
	}, err
}

type StringValue struct {
	basetypes.StringValue
}

// ValueStringPointer returns a pointer to the known string value, nil for a
// null string value, or a pointer to "" for an unknown string value.
func (b StringValue) ValueStringPointer() *string {
	if b.IsNull() {
		return nil
	}

	return bastionzero.PtrTo(b.ValueString())
}
