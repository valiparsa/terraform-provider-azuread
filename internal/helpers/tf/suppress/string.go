// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package suppress

import (
	"strings"

	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
)

func CaseDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return strings.EqualFold(old, new)
}
