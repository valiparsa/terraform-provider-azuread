// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package client

import (
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/approleassignedto"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/serviceprincipal"
	"github.com/valiparsa/terraform-provider-azuread/internal/common"
)

type Client struct {
	AppRoleAssignedToClient *approleassignedto.AppRoleAssignedToClient
	ServicePrincipalClient  *serviceprincipal.ServicePrincipalClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appRoleAssignedToClient, err := approleassignedto.NewAppRoleAssignedToClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(appRoleAssignedToClient.Client)

	servicePrincipalClient, err := serviceprincipal.NewServicePrincipalClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(servicePrincipalClient.Client)

	return &Client{
		AppRoleAssignedToClient: appRoleAssignedToClient,
		ServicePrincipalClient:  servicePrincipalClient,
	}, nil
}
