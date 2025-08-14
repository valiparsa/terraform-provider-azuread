// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package applications_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/applications/stable/application"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance/check"
	"github.com/valiparsa/terraform-provider-azuread/internal/clients"
	"github.com/valiparsa/terraform-provider-azuread/internal/services/applications/parse"
)

type ApplicationApiAccessResource struct{}

func TestAccApplicationApiAccess_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application_api_access", "test")
	r := ApplicationApiAccessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationApiAccess_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application_api_access", "test")
	data2 := acceptance.BuildTestData(t, "azuread_application_api_access", "test2")
	r := ApplicationApiAccessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("application_id").Exists(),
				check.That(data2.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
	})
}

func TestAccApplicationApiAccess_multipleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application_api_access", "test")
	data2 := acceptance.BuildTestData(t, "azuread_application_api_access", "test2")
	r := ApplicationApiAccessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("application_id").Exists(),
				check.That(data2.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("application_id").Exists(),
				check.That(data2.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("application_id").Exists(),
				check.That(data2.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
	})
}

func TestAccApplicationApiAccess_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application_api_access", "test")
	r := ApplicationApiAccessResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("api_client_id").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport(data)),
	})
}

func (r ApplicationApiAccessResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Applications.ApplicationClient

	id, err := parse.ParseApiAccessID(state.ID)
	if err != nil {
		return nil, err
	}

	applicationId := stable.NewApplicationID(id.ApplicationId)

	resp, err := client.GetApplication(ctx, applicationId, application.DefaultGetApplicationOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", applicationId, err)
	}

	app := resp.Model
	if app == nil {
		return nil, fmt.Errorf("retrieving %s: model was nil", applicationId)
	}

	if app.RequiredResourceAccess == nil {
		return pointer.To(false), nil
	}

	for _, api := range *app.RequiredResourceAccess {
		if strings.EqualFold(*api.ResourceAppId, id.ApiClientId) {
			return pointer.To(true), nil
		}
	}

	return pointer.To(false), nil
}

func (ApplicationApiAccessResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application_registration" "test" {
  display_name = "acctest-ApiAccess-%[1]d"
}

resource "azuread_application_api_access" "test" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0000-c000-000000000000"

  role_ids = [
    "9a5d68dd-52b0-4cc2-bd40-abcf44ac3a30",
    "dbb9058a-0e50-45d7-ae91-66909b5d4664",
  ]

  scope_ids = [
    "e1fe6dd8-ba31-4d61-89e7-88639da4683d",
  ]
}
`, data.RandomInteger, data.RandomPassword)
}

func (ApplicationApiAccessResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application_registration" "test" {
  display_name = "acctest-ApiAccess-%[1]d"
}

resource "azuread_application_api_access" "test" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0000-c000-000000000000"

  role_ids = [
    "9a5d68dd-52b0-4cc2-bd40-abcf44ac3a30",
    "dbb9058a-0e50-45d7-ae91-66909b5d4664",
  ]

  scope_ids = [
    "e1fe6dd8-ba31-4d61-89e7-88639da4683d",
  ]
}

resource "azuread_application_api_access" "test2" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0ff1-ce00-000000000000"

  role_ids = [
    "d13f72ca-a275-4b96-b789-48ebcc4da984",
  ]

  scope_ids = [
    "2beb830c-70d1-4f5b-a983-79cbdb0c6c6a",
  ]
}
`, data.RandomInteger, data.RandomPassword)
}

func (ApplicationApiAccessResource) multipleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application_registration" "test" {
  display_name = "acctest-ApiAccess-%[1]d"
}

resource "azuread_application_api_access" "test" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0000-c000-000000000000"

  scope_ids = [
    "e1fe6dd8-ba31-4d61-89e7-88639da4683d",
  ]
}

resource "azuread_application_api_access" "test2" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0ff1-ce00-000000000000"

  role_ids = [
    "d13f72ca-a275-4b96-b789-48ebcc4da984",
    "df021288-bdef-4463-88db-98f22de89214",
  ]

  scope_ids = [
    "2beb830c-70d1-4f5b-a983-79cbdb0c6c6a",
  ]
}
`, data.RandomInteger, data.RandomPassword)
}

func (ApplicationApiAccessResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application_registration" "test" {
  display_name = "acctest-ApiAccess-%[1]d"
}

resource "azuread_application_api_access" "test" {
  application_id = azuread_application_registration.test.id
  api_client_id  = "00000003-0000-0000-c000-000000000000"

  role_ids = [
    "9a5d68dd-52b0-4cc2-bd40-abcf44ac3a30",
    "dbb9058a-0e50-45d7-ae91-66909b5d4664",
  ]

  scope_ids = [
    "e1fe6dd8-ba31-4d61-89e7-88639da4683d",
  ]
}

resource "azuread_application_api_access" "import" {
  application_id = azuread_application_api_access.test.application_id
  api_client_id  = azuread_application_api_access.test.api_client_id
  role_ids       = azuread_application_api_access.test.role_ids
  scope_ids      = azuread_application_api_access.test.scope_ids
}
`, data.RandomInteger, data.RandomPassword)
}
