package tftypesgen

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
		).Draw(t, "SetWithValueOrNull")
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
		).Draw(t, "StringWithValueOrNull")
	})
}
