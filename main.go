// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14
// Modifications made on 2025-08-14

package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/valiparsa/terraform-provider-azuread/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        false,
		ProviderAddr: "registry.terraform.io/GPKbdZZb/forked-azuread",
		ProviderFunc: provider.AzureADProvider,
	}

	plugin.Serve(opts)
}
