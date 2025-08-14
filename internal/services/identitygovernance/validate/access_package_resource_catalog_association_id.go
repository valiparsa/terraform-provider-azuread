// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package validate

import (
	"github.com/valiparsa/terraform-provider-azuread/internal/services/identitygovernance/parse"
)

func AccessPackageResourceCatalogAssociationID(input string) (err error) {
	_, err = parse.AccessPackageResourceCatalogAssociationID(input)
	return
}
