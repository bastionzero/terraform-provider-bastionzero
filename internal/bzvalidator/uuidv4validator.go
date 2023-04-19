package bzvalidator

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = uuidV4Validator{}

// uuidV4Validator validates that a string Attribute's value is a valid UUID
// string
type uuidV4Validator struct{}

// Description describes the validation in plain text formatting.
func (validator uuidV4Validator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be a valid UUID v4 string")
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator uuidV4Validator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v uuidV4Validator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !isValidUUID(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// ValidUUIDV4 returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Is a valid UUID v4 string.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValidUUIDV4() validator.String {
	return uuidV4Validator{}
}

func isValidUUID(u string) bool {
	// Source: https://stackoverflow.com/a/46315070
	_, err := uuid.Parse(u)
	return err == nil
}
