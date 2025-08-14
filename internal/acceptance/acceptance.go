// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package acceptance

import (
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/valiparsa/terraform-provider-azuread/internal/provider"
)

var AzureADProvider *schema.Provider
var once sync.Once

func init() {
	if os.Getenv("TF_ACC") == "" {
		return
	}
	EnsureProvidersAreInitialised()
}

func EnsureProvidersAreInitialised() {
	once.Do(func() {
		AzureADProvider = provider.AzureADProvider()
	})
}
