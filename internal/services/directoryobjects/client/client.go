// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package client

import (
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/directoryobjects/stable/directoryobject"
	"github.com/valiparsa/terraform-provider-azuread/internal/common"
)

type Client struct {
	DirectoryObjectClient *directoryobject.DirectoryObjectClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	directoryObjectClient, err := directoryobject.NewDirectoryObjectClientWithBaseURI(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, err
	}
	o.Configure(directoryObjectClient.Client)

	return &Client{
		DirectoryObjectClient: directoryObjectClient,
	}, nil
}
