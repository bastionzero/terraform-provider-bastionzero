// Package bzgen provides property based test (PBT) generators for BastionZero
// API types
package bzgen

import (
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"pgregory.net/rapid"
)

func DatabaseAuthenticationConfigGen() *rapid.Generator[dbauthconfig.DatabaseAuthenticationConfig] {
	return rapid.Custom(func(t *rapid.T) dbauthconfig.DatabaseAuthenticationConfig {
		return dbauthconfig.DatabaseAuthenticationConfig{
			AuthenticationType:   rapid.Ptr(rapid.String(), true).Draw(t, "AuthenticationType"),
			CloudServiceProvider: rapid.Ptr(rapid.String(), true).Draw(t, "CloudServiceProvider"),
			Database:             rapid.Ptr(rapid.String(), true).Draw(t, "Database"),
			Label:                rapid.Ptr(rapid.String(), true).Draw(t, "Label"),
		}
	})
}
