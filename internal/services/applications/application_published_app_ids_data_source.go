// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package applications

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
)

func applicationPublishedAppIdsDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		ReadContext: func(_ context.Context, d *pluginsdk.ResourceData, _ interface{}) pluginsdk.Diagnostics {
			tf.Set(d, "result", environments.PublishedApis)
			d.SetId("appIds")
			return nil
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"result": {
				Description: "A mapping of application names and application IDs",
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}
