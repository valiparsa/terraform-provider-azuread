// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package provider

import (
	"github.com/valiparsa/terraform-provider-azuread/internal/sdk"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/administrativeunits"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/applications"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/approleassignments"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/conditionalaccess"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/directoryobjects"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/directoryroles"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/domains"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/groups"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/identitygovernance"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/invitations"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/policies"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/serviceprincipals"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/synchronization"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/userflows"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/users"
)

//go:generate go run ../tools/generator-services/main.go -path=../../

func SupportedTypedServices() []sdk.TypedServiceRegistration {
	return []sdk.TypedServiceRegistration{
		applications.Registration{},
		directoryroles.Registration{},
		domains.Registration{},
		policies.Registration{},
		identitygovernance.Registration{},
		serviceprincipals.Registration{},
	}
}

func SupportedUntypedServices() []sdk.UntypedServiceRegistration {
	return []sdk.UntypedServiceRegistration{
		administrativeunits.Registration{},
		applications.Registration{},
		approleassignments.Registration{},
		conditionalaccess.Registration{},
		directoryobjects.Registration{},
		directoryroles.Registration{},
		domains.Registration{},
		groups.Registration{},
		identitygovernance.Registration{},
		invitations.Registration{},
		policies.Registration{},
		serviceprincipals.Registration{},
		synchronization.Registration{},
		userflows.Registration{},
		users.Registration{},
	}
}
