package environment_test

import (
	"fmt"
	"strconv"
)

type environmentResourceTFConfigOptions struct {
	TFResourceName string

	Name                       *string
	Description                *string
	OfflineCleanupTimeoutHours *int
}

func surroundDoubleQuote(str string) string {
	return "\"" + str + "\""
}

func environmentResourceTFConfig(opts *environmentResourceTFConfigOptions) string {
	var name, description, cleanupTimeout string
	if opts.Name != nil {
		name = surroundDoubleQuote(*opts.Name)
	} else {
		name = "null"
	}
	if opts.Description != nil {
		description = surroundDoubleQuote(*opts.Description)
	} else {
		description = "null"
	}
	if opts.OfflineCleanupTimeoutHours != nil {
		cleanupTimeout = surroundDoubleQuote(strconv.Itoa(*opts.OfflineCleanupTimeoutHours))
	} else {
		cleanupTimeout = "null"
	}

	return fmt.Sprintf(`
resource "bastionzero_environment" "%s" {
  name   = %s
  description = %s
  offline_cleanup_timeout_hours = %s
}
`, opts.TFResourceName, name, description, cleanupTimeout)
}
