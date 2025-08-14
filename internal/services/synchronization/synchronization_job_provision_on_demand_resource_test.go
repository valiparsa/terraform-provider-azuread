// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package synchronization_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance"
	"github.com/valiparsa/terraform-provider-azuread/internal/clients"
)

type SynchronizationJobProvisionOnDemandResource struct{}

func TestAccSynchronizationJobProvisionOnDemand_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_synchronization_job_provision_on_demand", "test")
	r := SynchronizationJobProvisionOnDemandResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// The provisioned app isn't actually integrated so this will never work
			Config:      r.basic(data),
			ExpectError: regexp.MustCompile("CredentialsMissing"),
		},
	})
}

func (r SynchronizationJobProvisionOnDemandResource) Exists(_ context.Context, _ *clients.Client, _ *terraform.InstanceState) (*bool, error) {
	// Nothing to read
	return pointer.To(true), nil
}

func (SynchronizationJobProvisionOnDemandResource) template(data acceptance.TestData) string {
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

resource "azuread_synchronization_job" "test" {
  service_principal_id = data.azuread_service_principal.test.id
  template_id          = "dataBricks"
}

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
}
`, data.RandomInteger)
}

func (r SynchronizationJobProvisionOnDemandResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_synchronization_job_provision_on_demand" "test" {
  service_principal_id   = azuread_synchronization_job.test.service_principal_id
  synchronization_job_id = azuread_synchronization_job.test.id

  parameter {
    rule_id = "03f7d90d-bf71-41b1-bda6-aaf0ddbee5d8" // appears to be a global value

    subject {
      object_id        = azuread_group.test.id
      object_type_name = "Group"
    }
  }
}


`, r.template(data))
}
