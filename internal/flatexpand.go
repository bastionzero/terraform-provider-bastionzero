package internal

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Terraform Plugin Framework variants of standard flatteners and expanders.

// StringFromFramework converts a Framework String value to a string pointer. A
// null String is converted to a nil string pointer.
func StringFromFramework(_ context.Context, v types.String) *string {
	// Source: https://github.com/hashicorp/terraform-provider-aws/blob/0e19050852dadd4498d77467b8c2692b49881b22/internal/flex/framework.go
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return bastionzero.PtrTo(v.ValueString())
}

// ExpandFrameworkStringSet converts a framework Set value to a slice of string
// values.
//
// If the framework value is null or unknown, or if an error occurs when first
// converting the set into a slice of element type string, an empty slice is
// returned.
func ExpandFrameworkStringSet(ctx context.Context, tfSet types.Set) []string {
	return ExpandFrameworkSet(ctx, tfSet, func(m string) string {
		return m
	})
}

// ExpandFrameworkSet converts a framework Set value to a slice of values
// according to the specified map function f.
//
// If the framework value is null or unknown, or if an error occurs when first
// converting the set into a slice of type SetT, an empty slice is returned.
func ExpandFrameworkSet[SetT any, ExpandT any](ctx context.Context, set types.Set, f func(SetT) ExpandT) []ExpandT {
	if set.IsNull() || set.IsUnknown() {
		return []ExpandT{}
	}

	var data []SetT
	if set.ElementsAs(ctx, &data, false).HasError() {
		return []ExpandT{}
	}

	if len(data) == 0 {
		return []ExpandT{}
	}

	to := make([]ExpandT, len(data))
	for i, obj := range data {
		to[i] = f(obj)
	}

	return to
}

// FlattenFrameworkSet converts an arbitrary slice to a framework Set value of
// the specified elementType and map function f to do the conversion.
//
// If the slice has 0 elements, then an empty set is returned.
func FlattenFrameworkSet[T any](ctx context.Context, elementType attr.Type, list []T, f func(T) attr.Value) types.Set {
	elems := make([]attr.Value, len(list))
	for i, v := range list {
		elems[i] = f(v)
	}

	return types.SetValueMust(elementType, elems)
}
