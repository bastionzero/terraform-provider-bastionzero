// Package tfgen provides property based test (PBT) generators for BastionZero
// Terraform Provider schema model types
package tfgen

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"pgregory.net/rapid"
)

func SetWithValueOrNullGen(ctx context.Context, gen *rapid.Generator[basetypes.SetValue]) *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewSetNull(gen.Example(0).ElementType(ctx))),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "SetWithValueOrNull")
	})
}

func SetWithValueOrEmptyGen(ctx context.Context, gen *rapid.Generator[basetypes.SetValue]) *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewSetValueMust(gen.Example(0).ElementType(ctx), []attr.Value{})),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "SetWithValueOrEmpty")
	})
}

func SetWithValueOrNullOrEmptyGen(ctx context.Context, gen *rapid.Generator[basetypes.SetValue]) *rapid.Generator[basetypes.SetValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.SetValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewSetValueMust(gen.Example(0).ElementType(ctx), []attr.Value{})),
			rapid.Just(basetypes.NewSetNull(gen.Example(0).ElementType(ctx))),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "SetWithValueOrNullOrEmpty")
	})
}

func StringWithValueOrNullGen(gen *rapid.Generator[basetypes.StringValue]) *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewStringNull()),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "StringWithValueOrNull")
	})
}

func StringWithValueOrEmptyGen(gen *rapid.Generator[basetypes.StringValue]) *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewStringValue("")),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "StringWithValueOrEmpty")
	})
}

func StringWithValueOrNullOrEmptyGen(gen *rapid.Generator[basetypes.StringValue]) *rapid.Generator[basetypes.StringValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.StringValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewStringValue("")),
			rapid.Just(basetypes.NewStringNull()),
			rapid.Just(gen.Draw(t, "Value")),
		).Draw(t, "StringWithValueOrNullOrEmpty")
	})
}

func BoolWithValueOrNullGen(ctx context.Context) *rapid.Generator[basetypes.BoolValue] {
	return rapid.Custom(func(t *rapid.T) basetypes.BoolValue {
		return rapid.OneOf(
			rapid.Just(basetypes.NewBoolNull()),
			rapid.Just(basetypes.NewBoolValue(rapid.Bool().Draw(t, "Value"))),
		).Draw(t, "BoolWithValueOrNull")
	})
}

func Int64WithValueOrNullGen(ctx context.Context) *rapid.Generator[basetypes.Int64Value] {
	return rapid.Custom(func(t *rapid.T) basetypes.Int64Value {
		return rapid.OneOf(
			rapid.Just(basetypes.NewInt64Null()),
			rapid.Just(basetypes.NewInt64Value(rapid.Int64().Draw(t, "Value"))),
		).Draw(t, "Int64WithValueOrNull")
	})
}
