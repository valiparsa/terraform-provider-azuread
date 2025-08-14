// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package identitygovernance

import (
	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
	"github.com/valiparsa/terraform-provider-azuread/internal/sdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Identity Governance"
}

// AssociatedGitHubLabel is the issue/PR label which can be applied to PRs that include changes to this service package
func (r Registration) AssociatedGitHubLabel() string {
	return "feature/identity-governance"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Identity Governance",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azuread_access_package":              accessPackageDataSource(),
		"azuread_access_package_catalog":      accessPackageCatalogDataSource(),
		"azuread_access_package_catalog_role": accessPackageCatalogRoleDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azuread_access_package":                              accessPackageResource(),
		"azuread_access_package_assignment_policy":            accessPackageAssignmentPolicyResource(),
		"azuread_access_package_catalog":                      accessPackageCatalogResource(),
		"azuread_access_package_catalog_role_assignment":      accessPackageCatalogRoleAssignmentResource(),
		"azuread_access_package_resource_catalog_association": accessPackageResourceCatalogAssociationResource(),
		"azuread_access_package_resource_package_association": accessPackageResourcePackageAssociationResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		PrivilegedAccessGroupAssignmentScheduleResource{},
		PrivilegedAccessGroupEligibilityScheduleResource{},
	}
}
