// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package identitygovernance

import (
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/suppress"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/validation"
)

func schemaLocalizedContent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"default_text": {
				Description: "The default text of this question",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"localized_text": {
				Description: "The localized text of this question",
				Type:        pluginsdk.TypeList,
				Optional:    true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"language_code": {
							Description:  "The language code of this question content",
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.ISO639Language,
						},

						"content": {
							Description: "The localized content of this question",
							Type:        pluginsdk.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func schemaUserSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"subject_type": {
				Description:      "Type of users",
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"ConnectedOrganizationMembers",
					"ExternalSponsors",
					"GroupMembers",
					"InternalSponsors",
					"RequestorManager",
					"SingleUser",
					"TargetUserSponsors",
				}, true),
			},

			"backup": {
				Description: "For a user in an approval stage, this property indicates whether the user is a backup fallback approver",
				Type:        pluginsdk.TypeBool,
				Optional:    true,
			},

			"object_id": {
				Description: "The object ID of the subject",
				Type:        pluginsdk.TypeString,
				Optional:    true,
			},
		},
	}
}
