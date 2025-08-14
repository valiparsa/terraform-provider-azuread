// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package migrations

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
)

func ResourceDirectoryRoleAssignmentInstanceResourceV0() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"role_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"principal_object_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"app_scope_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"directory_scope_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceDirectoryRoleAssignmentInstanceStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Println("[DEBUG] Migrating ID from v0 to v1 format")
	oldId := rawState["id"].(string)
	newId := stable.NewRoleManagementDirectoryRoleAssignmentID(oldId)
	rawState["id"] = newId.ID()
	return rawState, nil
}
