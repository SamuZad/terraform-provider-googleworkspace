// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	directory "google.golang.org/api/admin/directory/v1"
)

// groupMemberRoles is a group's membership as the Directory API lists it,
// keyed by lower-cased email (the API lower-cases addresses) with the
// member's role as the value.
type groupMemberRoles map[string]string

// expectedGroupMemberRoles is the membership a resource's `members` set
// declares, in the shape listGroupMemberRoles returns.
func expectedGroupMemberRoles(members []interface{}) groupMemberRoles {
	expected := make(groupMemberRoles, len(members))
	for _, raw := range members {
		member := raw.(map[string]interface{})
		expected[strings.ToLower(member["email"].(string))] = member["role"].(string)
	}
	return expected
}

// listGroupMemberRoles lists the group's direct members once.
func listGroupMemberRoles(ctx context.Context, membersService *directory.MembersService, groupId string) (groupMemberRoles, error) {
	listed := groupMemberRoles{}
	err := retryNotFound(ctx, func() error {
		for k := range listed {
			delete(listed, k)
		}
		return membersService.List(groupId).MaxResults(200).Pages(ctx, func(resp *directory.Members) error {
			for _, member := range resp.Members {
				listed[strings.ToLower(member.Email)] = member.Role
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return listed, nil
}

// waitForGroupMembers polls `list` until it returns exactly `expected`, or
// `timeout` passes.
//
// Member inserts, updates and deletes are eventually consistent: a List
// issued right after them can still omit a member just inserted, carry one
// just deleted, or show the old role. Reading state from that snapshot
// records a membership the configuration never declared, and the next plan
// then re-inserts (409) or re-deletes (404) members that are already right.
// Waiting for the listing to match what was written is the same contract
// the group and group_member resources keep through consistencyCheck.
func waitForGroupMembers(ctx context.Context, timeout time.Duration, expected groupMemberRoles, list func() (groupMemberRoles, error)) error {
	return retryTimeDuration(ctx, timeout, func() error {
		listed, err := list()
		if err != nil {
			return fmt.Errorf("unexpected error while waiting for group members: %s", err)
		}
		if groupMembersMatch(expected, listed) {
			return nil
		}
		// The "timed out while waiting" wording is what IsNotConsistent keys on
		// to make retryTimeDuration try again.
		return fmt.Errorf("timed out while waiting for group members to be consistent: expected %d members, listed %d", len(expected), len(listed))
	})
}

func groupMembersMatch(expected, listed groupMemberRoles) bool {
	if len(expected) != len(listed) {
		return false
	}
	for email, role := range expected {
		if listedRole, ok := listed[email]; !ok || listedRole != role {
			return false
		}
	}
	return true
}
