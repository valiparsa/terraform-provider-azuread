// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package synchronization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/synchronizationjob"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance/check"
	"github.com/valiparsa/terraform-provider-azuread/internal/clients"
)

type SynchronizationJobResource struct{}

func TestAccSynchronizationJob(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"synchronizationJob": {
			"basic":    testAccSynchronizationJob_basic,
			"disabled": testAccSynchronizationJob_disabled,
		},
	})
}

func testAccSynchronizationJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_synchronization_job", "test")
	r := SynchronizationJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("template_id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccSynchronizationJob_disabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_synchronization_job", "test")
	r := SynchronizationJobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("template_id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (r SynchronizationJobResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.ServicePrincipals.SynchronizationJobClient

	id, err := stable.ParseServicePrincipalIdSynchronizationJobID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing synchronization job ID: %v", err)
	}

	resp, err := client.GetSynchronizationJob(ctx, *id, synchronizationjob.DefaultGetSynchronizationJobOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s", id)
	}

	return pointer.To(true), nil
}

func (SynchronizationJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_client_config" "test" {}

data "azuread_application_template" "test" {
  display_name = "Azure Databricks SCIM Provisioning Connector"
}

resource "azuread_application_from_template" "test" {
  display_name = "acctestSynchronizationJob-%[1]d"
  template_id  = data.azuread_application_template.test.template_id
}

data "azuread_service_principal" "test" {
  object_id = azuread_application_from_template.test.service_principal_object_id
}
`, data.RandomInteger)
}

func (r SynchronizationJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_synchronization_job" "test" {
  service_principal_id = data.azuread_service_principal.test.id
  template_id          = "dataBricks"
}
`, r.template(data))
}

func (r SynchronizationJobResource) disabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_synchronization_job" "test" {
  service_principal_id = data.azuread_service_principal.test.id
  template_id          = "dataBricks"
  enabled              = false
}
`, r.template(data))
}
