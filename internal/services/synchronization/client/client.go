// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package client

import (
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/serviceprincipal"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/synchronizationjob"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/serviceprincipals/stable/synchronizationsecret"
	"github.com/valiparsa/terraform-provider-azuread/internal/common"
)

type Client struct {
	ServicePrincipalClient      *serviceprincipal.ServicePrincipalClient
	SynchronizationJobClient    *synchronizationjob.SynchronizationJobClient
	SynchronizationSecretClient *synchronizationsecret.SynchronizationSecretClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	servicePrincipalClient, err := serviceprincipal.NewServicePrincipalClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(servicePrincipalClient.Client)

	synchronizationJobClient, err := synchronizationjob.NewSynchronizationJobClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(synchronizationJobClient.Client)

	synchronizationSecretClient, err := synchronizationsecret.NewSynchronizationSecretClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(synchronizationSecretClient.Client)

	return &Client{
		ServicePrincipalClient:      servicePrincipalClient,
		SynchronizationJobClient:    synchronizationJobClient,
		SynchronizationSecretClient: synchronizationSecretClient,
	}, nil
}
