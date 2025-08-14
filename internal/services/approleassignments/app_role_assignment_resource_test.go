// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package approleassignments_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/approleassignedto"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance/check"
	"github.com/valiparsa/terraform-provider-azuread/internal/clients"
)

type AppRoleAssignmentResource struct{}

func TestAccAppRoleAssignment_servicePrincipalForMsGraph(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_app_role_assignment", "test")
	r := AppRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipalForMsGraph(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppRoleAssignment_servicePrincipalForTenantApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_app_role_assignment", "test_admin")
	r := AppRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipalForTenantApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azuread_app_role_assignment.test_query").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppRoleAssignment_groupForTenantApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_app_role_assignment", "test")
	r := AppRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.groupForTenantApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppRoleAssignment_groupForTenantAppWithoutRole(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_app_role_assignment", "test")
	r := AppRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.groupForTenantAppWithoutRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAppRoleAssignment_userForTenantApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_app_role_assignment", "test")
	r := AppRoleAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userForTenantApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AppRoleAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.AppRoleAssignments.AppRoleAssignedToClient

	id, err := stable.ParseServicePrincipalIdAppRoleAssignedToID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing App Role Assignment ID: %v", err)
	}

	resp, err := client.GetAppRoleAssignedTo(ctx, *id, approleassignedto.DefaultGetAppRoleAssignedToOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("failed to retrieve %s: %+v", id, err)
	}

	appRoleAssignment := resp.Model
	if appRoleAssignment == nil {
		return nil, fmt.Errorf("retrieving %s: model was nil", id)
	}

	return pointer.To(true), nil
}

func (AppRoleAssignmentResource) servicePrincipalForMsGraph(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_service_principal" "msgraph" {
  client_id    = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
  use_existing = true
}

resource "azuread_application" "test" {
  display_name = "acctest-appRoleAssignment-%[1]d"

  required_resource_access {
    resource_app_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph

    resource_access {
      id   = azuread_service_principal.msgraph.app_role_ids["User.Read.All"]
      type = "Role"
    }

    resource_access {
      id   = azuread_service_principal.msgraph.oauth2_permission_scope_ids["User.ReadWrite"]
      type = "Scope"
    }
  }
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application.test.client_id
}

resource "azuread_app_role_assignment" "test" {
  app_role_id         = azuread_service_principal.msgraph.app_role_ids["User.Read.All"]
  principal_object_id = azuread_service_principal.test.object_id
  resource_object_id  = azuread_service_principal.msgraph.object_id
}
`, data.RandomInteger)
}

func (AppRoleAssignmentResource) tenantAppTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "internal" {
  display_name = "acctest-AppRoleAssignment-internal-%[1]d"

  app_role {
    allowed_member_types = ["Application", "User"]
    description          = "Admins can perform all task actions"
    display_name         = "Admin"
    enabled              = true
    id                   = "%[2]s"
    value                = "Admin.All"
  }

  app_role {
    allowed_member_types = ["Application"]
    description          = "Apps can query the database"
    display_name         = "Query"
    enabled              = true
    id                   = "%[3]s"
    value                = "Query.All"
  }
}

resource "azuread_service_principal" "internal" {
  client_id = azuread_application.internal.client_id
}
`, data.RandomInteger, data.UUID(), data.UUID())
}

func (r AppRoleAssignmentResource) servicePrincipalForTenantApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test" {
  display_name = "acctest-appRoleAssignment-%[2]d"

  required_resource_access {
    resource_app_id = azuread_application.internal.client_id

    resource_access {
      id   = azuread_service_principal.internal.app_role_ids["Admin.All"]
      type = "Role"
    }

    resource_access {
      id   = azuread_service_principal.internal.app_role_ids["Query.All"]
      type = "Role"
    }
  }
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application.test.client_id
}

resource "azuread_app_role_assignment" "test_admin" {
  app_role_id         = azuread_service_principal.internal.app_role_ids["Admin.All"]
  principal_object_id = azuread_service_principal.test.object_id
  resource_object_id  = azuread_service_principal.internal.object_id
}

resource "azuread_app_role_assignment" "test_query" {
  app_role_id         = azuread_service_principal.internal.app_role_ids["Query.All"]
  principal_object_id = azuread_service_principal.test.object_id
  resource_object_id  = azuread_service_principal.internal.object_id
}
`, r.tenantAppTemplate(data), data.RandomInteger)
}

func (r AppRoleAssignmentResource) groupForTenantApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctest-appRoleAssignment-%[2]d"
  security_enabled = true
}

resource "azuread_app_role_assignment" "test" {
  app_role_id         = azuread_service_principal.internal.app_role_ids["Admin.All"]
  principal_object_id = azuread_group.test.object_id
  resource_object_id  = azuread_service_principal.internal.object_id
}
`, r.tenantAppTemplate(data), data.RandomInteger)
}

func (r AppRoleAssignmentResource) groupForTenantAppWithoutRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "internal" {
  display_name = "acctest-AppRoleAssignment-internal-%[1]d"
}

resource "azuread_service_principal" "internal" {
  client_id = azuread_application.internal.client_id
}

resource "azuread_group" "test" {
  display_name     = "acctest-appRoleAssignment-%[1]d"
  security_enabled = true
}

resource "azuread_app_role_assignment" "test" {
  app_role_id         = "00000000-0000-0000-0000-000000000000"
  principal_object_id = azuread_group.test.object_id
  resource_object_id  = azuread_service_principal.internal.object_id
}
`, data.RandomInteger)
}

func (r AppRoleAssignmentResource) userForTenantApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  display_name        = "acctest-appRoleAssignment-%[2]d"
  password            = "%[3]s"
  user_principal_name = "acctest-AppRoleAssignment-%[2]d@${data.azuread_domains.test.domains.0.domain_name}"
}

resource "azuread_app_role_assignment" "test" {
  app_role_id         = azuread_service_principal.internal.app_role_ids["Admin.All"]
  principal_object_id = azuread_user.test.object_id
  resource_object_id  = azuread_service_principal.internal.object_id
}
`, r.tenantAppTemplate(data), data.RandomInteger, data.RandomPassword)
}
