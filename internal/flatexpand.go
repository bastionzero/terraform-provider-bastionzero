package internal

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/typesext"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Terraform Plugin Framework variants of standard flatteners and expanders.
// Source: https://github.com/hashicorp/terraform-provider-aws/blob/0e19050852dadd4498d77467b8c2692b49881b22/internal/flex/framework.go

// FlattenFrameworkStringList converts a slice of string pointers to a framework
// List value.
//
// A nil slice is converted to a null List. An empty slice is converted to a
// null List.
func FlattenFrameworkStringList(_ context.Context, vs []*string) types.List {
	if len(vs) == 0 {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = typesext.StringPointerValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
}

// ExpandFrameworkStringList converts a framework List value to a slice of
// string pointers.
//
// If the framework value is null or unknown, a nil slice is returned.
func ExpandFrameworkStringList(ctx context.Context, list types.List) []*string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var vl []*string

	if list.ElementsAs(ctx, &vl, false).HasError() {
		return nil
	}

	return vl
}

// FlattenFrameworkStringValueList converts a slice of string values to a
// framework List value.
//
// A nil slice is converted to a null List. An empty slice is converted to a
// null List.
func FlattenFrameworkStringValueList(_ context.Context, vs []string) types.List {
	if len(vs) == 0 {
		return types.ListNull(types.StringType)
	}

	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, elems)
}

// ExpandFrameworkStringValueList converts a framework List value to a slice of
// string values.
//
// If the framework value is null or unknown, a nil slice is returned.
func ExpandFrameworkStringValueList(ctx context.Context, list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var vl []string

	if list.ElementsAs(ctx, &vl, false).HasError() {
		return nil
	}

	return vl
}

// FlattenFrameworkStringValueMap converts a map of string values to a framework
// Map value.
//
// A nil map is converted to an empty (non-null) Map.
func FlattenFrameworkStringValueMap(_ context.Context, m map[string]string) types.Map {
	elems := make(map[string]attr.Value, len(m))

	for k, v := range m {
		elems[k] = types.StringValue(v)
	}

	return types.MapValueMust(types.StringType, elems)
}

// ExpandFrameworkStringValueMap converts a framework Map value to a map of
// string values.
//
// If the framework value is null or unknown, a nil map is returned.
func ExpandFrameworkStringValueMap(ctx context.Context, set types.Map) map[string]string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var m map[string]string

	if set.ElementsAs(ctx, &m, false).HasError() {
		return nil
	}

	return m
}

// FlattenFrameworkStringSet converts a slice of string pointers to a framework
// Set value.
//
// A nil slice is converted to a null Set. An empty slice is converted to a null
// Set.
func FlattenFrameworkStringSet(_ context.Context, vs []*string) types.Set {
	if len(vs) == 0 {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = typesext.StringPointerValue(v)
	}

	return types.SetValueMust(types.StringType, elems)
}

// ExpandFrameworkStringValueMap converts a framework Set value to a slice of
// string pointers.
//
// If the framework value is null or unknown, a nil slice is returned.
func ExpandFrameworkStringSet(ctx context.Context, set types.Set) []*string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var vs []*string

	if set.ElementsAs(ctx, &vs, false).HasError() {
		return nil
	}

	return vs
}

// FlattenFrameworkStringValueSet converts a slice of string values to a
// framework Set value.
//
// A nil slice is converted to a null Set. An empty slice is converted to a null
// Set.
func FlattenFrameworkStringValueSet(_ context.Context, vs []string) types.Set {
	if len(vs) == 0 {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, len(vs))

	for i, v := range vs {
		elems[i] = types.StringValue(v)
	}

	return types.SetValueMust(types.StringType, elems)
}

// ExpandFrameworkStringValueSet converts a framework Set value to a slice of
// string values.
//
// If the framework value is null or unknown, a nil slice is returned.
func ExpandFrameworkStringValueSet(ctx context.Context, set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var vs []string

	if set.ElementsAs(ctx, &vs, false).HasError() {
		return nil
	}

	return vs
}

// StringFromFramework converts a Framework String value to a string pointer. A
// null String is converted to a nil string pointer.
func StringFromFramework(_ context.Context, v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return bastionzero.PtrTo(v.ValueString())
}

// StringValueToFramework converts a string value to a Framework String value.
// An empty string is converted to a null String.
func StringValueToFramework[T ~string](_ context.Context, v T) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(string(v))
}

// StringValueToFrameworkLegacy converts a string value to a Framework String
// value. An empty string is left as an empty String.
func StringValueToFrameworkLegacy[T ~string](_ context.Context, v T) types.String {
	return types.StringValue(string(v))
}
