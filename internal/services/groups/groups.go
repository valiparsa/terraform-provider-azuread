// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package groups

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/beta"
	groupBeta "github.com/hashicorp/go-azure-sdk/microsoft-graph/groups/beta/group"
	memberBeta "github.com/hashicorp/go-azure-sdk/microsoft-graph/groups/beta/member"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// groupCreationWriteRetryFunc retries a write operation against a freshly created group when it fails due to Azure AD
// replication lag. Two eventual-consistency conditions are tolerated:
//
//   - 404 Not Found: the group itself has not yet replicated to the replica serving the request.
//   - 403 Authorization_RequestDenied: the group has replicated, but the calling principal's authorization context
//     over it (e.g. the ownership that grants write access) has not, so the replica rejects the write as if the
//     caller had insufficient privileges.
//
// Returning true retries the request; the retry is bounded by the calling operation's context timeout (the resource
// Create timeout), so replication is given time to catch up rather than failing the apply outright.
func groupCreationWriteRetryFunc(resp *http.Response, o *odata.OData) (bool, error) {
	if response.WasNotFound(resp) {
		return true, nil
	}
	if response.WasForbidden(resp) && o != nil && o.Error != nil {
		return o.Error.Match("Authorization_RequestDenied"), nil
	}
	return false, nil
}

func groupDefaultMailNickname() string {
	charSet := "0123456789abcdef"
	result := make([]byte, 9)
	for i := 0; i < 9; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	resultString := string(result)
	return resultString[:8] + "-" + resultString[8:]
}

func groupFindByName(ctx context.Context, client *groupBeta.GroupClient, displayName string) (*[]beta.Group, error) {
	options := groupBeta.ListGroupsOperationOptions{
		Filter: pointer.To(fmt.Sprintf("displayName eq '%s'", displayName)),
	}

	resp, err := client.ListGroups(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("unable to list Groups with filter %q: %v", *options.Filter, err)
	}

	result := make([]beta.Group, 0)
	if resp.Model != nil {
		for _, group := range *resp.Model {
			if group.DisplayName != nil && group.DisplayName.GetOrZero() == displayName {
				result = append(result, group)
			}
		}
	}

	return &result, nil
}

func groupGetAdditional(ctx context.Context, client *groupBeta.GroupClient, id beta.GroupId) (*beta.Group, error) {
	options := groupBeta.GetGroupOperationOptions{
		Select: &[]string{
			"allowExternalSenders",
			"autoSubscribeNewMembers",
			"hideFromAddressLists",
			"hideFromOutlookClients",
		},
	}

	resp, err := client.GetGroup(ctx, id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			// API returns 404 when these M365-only fields are requested for a group in a non-M365 tenant, so we
			// don't raise an error in this case and proceed as if they are not set.
			// See https://github.com/microsoftgraph/msgraph-metadata/issues/333
			return nil, nil
		}

		return nil, fmt.Errorf("retrieving additional fields: %+v", err)
	}

	return resp.Model, nil
}

func groupGetMember(ctx context.Context, client *memberBeta.MemberClient, id beta.GroupIdMemberId) (*beta.DirectoryObject, error) {
	options := memberBeta.ListMembersOperationOptions{
		Filter: pointer.To(fmt.Sprintf("id eq '%s'", id.DirectoryObjectId)),
	}

	resp, err := client.ListMembers(ctx, beta.NewGroupID(id.GroupId), options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if resp.Model != nil {
		for _, member := range *resp.Model {
			if member.DirectoryObject().Id != nil && *member.DirectoryObject().Id == id.DirectoryObjectId {
				return &member, nil
			}
		}
	}

	return nil, nil
}
