// bztftest provides common testing functionality used throughout unit tests and
// acceptance tests for the bastionzero provider
package bztftest

import (
	bzapi "github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// ProviderConfig is a shared configuration to combine with the actual test
	// configuration so the BastionZero client is properly configured.
	ProviderConfig = `
provider "bastionzero" {}
`
)

var (
	TestProvider provider.Provider
	// TestProtoV6ProviderFactories are used to instantiate a provider during
	// testing. The factory function will be invoked for every Terraform CLI
	// command executed to create a provider server to which the CLI can
	// reattach.
	TestProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
)

func init() {
	TestProvider = bastionzero.New("test")()
	TestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"bastionzero": providerserver.NewProtocol6WithError(TestProvider),
	}

}

// GetBastionZeroClient attempts to load the BastionZero API client that was
// created during provider configuration. Panics if the client is not set.
func GetBastionZeroClient() *bzapi.Client {
	client := TestProvider.(*bastionzero.BastionZeroProvider).Client
	if client == nil {
		panic("Provider's Client is nil. Please ensure Configure() is called")
	}

	return client
}
