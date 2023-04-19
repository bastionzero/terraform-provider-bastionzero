package bastionzero_test

import (
	"regexp"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProviderConfig_InvalidAPISecret(t *testing.T) {
	// Clear env-var if set
	closer := acctest.SetEnvironmentVariables(map[string]string{
		"BASTIONZERO_API_SECRET": "",
	})
	t.Cleanup(closer)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Missing secret
				Config: `
					provider "bastionzero" {}
					resource "bastionzero_environment" "test" {
						name = "test"
					}
				`,
				ExpectError: regexp.MustCompile(`Missing BastionZero API Secret`),
			},
			{
				// Bad secret
				Config: `
					provider "bastionzero" {
						api_secret = "foo"
					}
					resource "bastionzero_environment" "test" {
						name = "test"
					}
				`,
				ExpectError: regexp.MustCompile(`Unable to create BastionZero API Client`),
			},
		},
	})
}
