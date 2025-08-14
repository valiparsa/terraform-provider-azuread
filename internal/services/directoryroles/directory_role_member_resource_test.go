// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package directoryroles_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/directoryroles/stable/member"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance"
	"github.com/valiparsa/terraform-provider-azuread/internal/acceptance/check"
	"github.com/valiparsa/terraform-provider-azuread/internal/clients"
)

type DirectoryRoleMemberResource struct{}

func TestAccDirectoryRoleMember_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_directory_role_member", "test")
	r := DirectoryRoleMemberResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_object_id").IsUuid(),
				check.That(data.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDirectoryRoleMember_user(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_directory_role_member", "testA")
	r := DirectoryRoleMemberResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneUser(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_object_id").IsUuid(),
				check.That(data.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDirectoryRoleMember_multipleUser(t *testing.T) {
	dataA := acceptance.BuildTestData(t, "azuread_directory_role_member", "testA")
	dataB := acceptance.BuildTestData(t, "azuread_directory_role_member", "testB")
	r := DirectoryRoleMemberResource{}

	dataA.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.oneUser(dataA),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("role_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		dataA.ImportStep(),
		{
			Config: r.twoUsers(dataA),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("role_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
				check.That(dataB.ResourceName).ExistsInAzure(r),
				check.That(dataB.ResourceName).Key("role_object_id").IsUuid(),
				check.That(dataB.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		dataA.ImportStep(),
		dataB.ImportStep(),
		{
			Config: r.oneUser(dataA),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("role_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		dataA.ImportStep(),
	})
}

func TestAccDirectoryRoleMember_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_directory_role_member", "test")
	r := DirectoryRoleMemberResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.servicePrincipal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport(data)),
	})
}

func (r DirectoryRoleMemberResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.DirectoryRoles.DirectoryRoleMemberClient

	id, err := stable.ParseDirectoryRoleIdMemberID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing Directory Role Member ID: %v", err)
	}

	options := member.ListMembersOperationOptions{
		Filter: pointer.To(fmt.Sprintf("id eq '%s'", id.DirectoryObjectId)),
	}
	resp, err := client.ListMembers(ctx, stable.NewDirectoryRoleID(id.DirectoryRoleId), options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("failed to retrieve %s: %v", id, err)
	}

	if resp.Model != nil {
		for _, member := range *resp.Model {
			if pointer.From(member.DirectoryObject().Id) == id.DirectoryObjectId {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}

func (DirectoryRoleMemberResource) templateThreeUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "testA" {
  user_principal_name = "acctestUser.%[1]d.A@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-A"
  password            = "%[2]s"
}

resource "azuread_user" "testB" {
  user_principal_name = "acctestUser.%[1]d.B@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-B"
  mail_nickname       = "acctestUser-%[1]d-B"
  password            = "%[2]s"
}

resource "azuread_user" "testC" {
  user_principal_name = "acctestUser.%[1]d.C@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-C"
  password            = "%[2]s"
}
`, data.RandomInteger, data.RandomPassword)
}

func (r DirectoryRoleMemberResource) servicePrincipal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test" {
  display_name = "acctestServicePrincipal-%[2]d"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application.test.client_id
}

resource "azuread_directory_role_member" "test" {
  role_object_id   = azuread_directory_role.test.object_id
  member_object_id = azuread_service_principal.test.object_id
}
`, DirectoryRoleResource{}.byTemplateId(data), data.RandomInteger)
}

func (r DirectoryRoleMemberResource) oneUser(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azuread_directory_role_member" "testA" {
  role_object_id   = azuread_directory_role.test.object_id
  member_object_id = azuread_user.testA.object_id
}
`, DirectoryRoleResource{}.byTemplateId(data), r.templateThreeUsers(data))
}

func (r DirectoryRoleMemberResource) twoUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azuread_directory_role_member" "testA" {
  role_object_id   = azuread_directory_role.test.object_id
  member_object_id = azuread_user.testA.object_id
}

resource "azuread_directory_role_member" "testB" {
  role_object_id   = azuread_directory_role.test.object_id
  member_object_id = azuread_user.testB.object_id
}
`, DirectoryRoleResource{}.byTemplateId(data), r.templateThreeUsers(data))
}

func (r DirectoryRoleMemberResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_directory_role_member" "import" {
  role_object_id   = azuread_directory_role_member.test.role_object_id
  member_object_id = azuread_directory_role_member.test.member_object_id
}
`, r.servicePrincipal(data))
}
