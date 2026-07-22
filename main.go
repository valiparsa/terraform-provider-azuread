// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14
// Modifications made on 2025-08-14

package main

import (
	"flag"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/valiparsa/terraform-provider-azuread/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// The provider's registry address doubles as the reattach key that Terraform
	// matches against a configuration's required_providers "source". Allow it to be
	// overridden (e.g. to "registry.terraform.io/hashicorp/azuread") for local
	// debugging against configs that reference the upstream source.
	providerAddr := "registry.terraform.io/GPKbdZZb/forked-azuread"
	if v := os.Getenv("PROVIDER_ADDR"); v != "" {
		providerAddr = v
	}

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: providerAddr,
		ProviderFunc: provider.AzureADProvider,
	}

	plugin.Serve(opts)
}
