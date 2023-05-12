package main

import (
	"context"
	"flag"
	"log"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for
// the registry/website

// Run terraform format on examples to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// Run the genlinks script to find links to the BastionZero docs website in the provider docs
//go:generate bash scripts/genlinks.sh

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

// ProviderAddr contains the full name for this terraform provider.
const ProviderAddr = "registry.terraform.io/bastionzero/bastionzero"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: ProviderAddr,
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), bastionzero.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
