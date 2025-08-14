// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package identitygovernance

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/go-azure-sdk/sdk/nullable"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/validation"
	"github.com/valiparsa/terraform-provider-azuread/internal/sdk"
)

type PrivilegedAccessGroupScheduleModel struct {
	AssignmentType      string `tfschema:"assignment_type"`
	Duration            string `tfschema:"duration"`
	ExpirationDate      string `tfschema:"expiration_date"`
	GroupId             string `tfschema:"group_id"`
	Justification       string `tfschema:"justification"`
	PermanentAssignment bool   `tfschema:"permanent_assignment"`
	PrincipalId         string `tfschema:"principal_id"`
	StartDate           string `tfschema:"start_date"`
	Status              string `tfschema:"status"`
	TicketNumber        string `tfschema:"ticket_number"`
	TicketSystem        string `tfschema:"ticket_system"`
}

func privilegedAccessGroupScheduleArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"group_id": {
			Description:      "The ID of the Group representing the scope of the assignment",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ValidateDiag(validation.IsUUID),
		},

		"principal_id": {
			Description:      "The ID of the Principal assigned to the schedule",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ValidateDiag(validation.IsUUID),
		},

		"assignment_type": {
			Description:      "The ID of the assignment to the group",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ValidateDiag(validation.StringInSlice(stable.PossibleValuesForPrivilegedAccessGroupRelationships(), false)),
		},

		"start_date": {
			Description:           "The date that this assignment starts, formatted as an RFC3339 date string in UTC (e.g. 2018-01-01T01:02:03Z)",
			Type:                  pluginsdk.TypeString,
			Optional:              true,
			Computed:              true,
			ValidateDiagFunc:      validation.ValidateDiag(validation.IsRFC3339Time),
			DiffSuppressOnRefresh: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				// Suppress diffs if the start date is in the past
				oldTime, err := time.Parse(time.RFC3339, old)
				if err == nil {
					return oldTime.Before(time.Now())
				}
				// Suppress diffs if the new date is within 5 minutes of the old date
				// Activation of a future start time is never exactly at the requested time
				newTime, err := time.Parse(time.RFC3339, new)
				if err == nil {
					return newTime.Before(oldTime.Add(5 * time.Minute))
				}
				return false
			},
		},

		"expiration_date": {
			Description:      "The date that this assignment expires, formatted as an RFC3339 date string in UTC (e.g. 2018-01-01T01:02:03Z)",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ConflictsWith:    []string{"duration"},
			ValidateDiagFunc: validation.ValidateDiag(validation.IsRFC3339Time),
		},

		"duration": {
			Description:      "The duration of the assignment, formatted as an ISO8601 duration string (e.g. P3D for 3 days)",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"expiration_date"},
			ValidateDiagFunc: validation.ValidateDiag(validation.StringIsNotEmpty),
		},

		"permanent_assignment": {
			Description: "Is the assignment permanent",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
		},

		"justification": {
			Description:      "The justification for the assignment",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ValidateDiag(validation.StringIsNotEmpty),
		},

		"ticket_number": {
			Description:      "The ticket number authorising the assignment",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			RequiredWith:     []string{"ticket_system"},
			ValidateDiagFunc: validation.ValidateDiag(validation.StringIsNotEmpty),
		},

		"ticket_system": {
			Description:      "The ticket system authorising the assignment",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			RequiredWith:     []string{"ticket_number"},
			ValidateDiagFunc: validation.ValidateDiag(validation.StringIsNotEmpty),
		},
	}
}

func privilegedAccessGroupScheduleAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"status": {
			Description: "The status of the schedule",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func buildScheduleRequest(model *PrivilegedAccessGroupScheduleModel, metadata *sdk.ResourceMetaData) (*stable.RequestSchedule, error) {
	schedule := stable.RequestSchedule{
		Expiration:    &stable.ExpirationPattern{},
		StartDateTime: nullable.NoZero(model.StartDate),
	}
	var startDate, expiryDate time.Time

	if model.StartDate != "" {
		var err error
		startDate, err = time.Parse(time.RFC3339, model.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %+v", model.StartDate, err)
		}
	}

	switch {
	case model.ExpirationDate != "":
		var err error
		expiryDate, err = time.Parse(time.RFC3339, model.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %+v", model.ExpirationDate, err)
		}

		if model.StartDate != "" && expiryDate.Before(startDate.Add(5*time.Minute)) {
			return nil, fmt.Errorf("`expiration_date` must be at least 5 minutes after `start_date`")
		}

		if metadata.ResourceData.HasChange("expiry_date") && expiryDate.Before(time.Now().Add(5*time.Minute)) {
			return nil, fmt.Errorf("`expiration_date` must be at least 5 minutes in the future")
		}

		schedule.Expiration.EndDateTime = nullable.Value(model.ExpirationDate)
		schedule.Expiration.Type = pointer.To(stable.ExpirationPatternType_AfterDateTime)

	case model.Duration != "":
		schedule.Expiration.Duration = nullable.Value(model.Duration)
		schedule.Expiration.Type = pointer.To(stable.ExpirationPatternType_AfterDuration)

	case model.PermanentAssignment:
		schedule.Expiration.Type = pointer.To(stable.ExpirationPatternType_NoExpiration)

	default:
		return nil, fmt.Errorf("either `expiration_date` or `duration` must be set, or `permanent_assignment` must be true")
	}

	return &schedule, nil
}
