// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitForGroupMembers_convergesOnceTheListingCatchesUp(t *testing.T) {
	t.Parallel()

	expected := groupMemberRoles{"a@example.com": "OWNER", "b@example.com": "MEMBER"}
	snapshots := []groupMemberRoles{
		{"a@example.com": "OWNER"},                           // the insert of b@ not yet visible
		{"a@example.com": "OWNER", "b@example.com": "OWNER"}, // b@'s role update not yet visible
		{"a@example.com": "OWNER", "b@example.com": "MEMBER"},
	}
	calls := 0
	list := func() (groupMemberRoles, error) {
		snapshot := snapshots[calls]
		if calls < len(snapshots)-1 {
			calls++
		}
		return snapshot, nil
	}

	if err := waitForGroupMembers(context.Background(), 20*time.Second, expected, list); err != nil {
		t.Fatalf("expected convergence, got %s", err)
	}
	if calls != len(snapshots)-1 {
		t.Fatalf("expected %d stale listings before convergence, got %d", len(snapshots)-1, calls)
	}
}

func TestWaitForGroupMembers_timesOutWhileStale(t *testing.T) {
	t.Parallel()

	expected := groupMemberRoles{"a@example.com": "OWNER", "b@example.com": "MEMBER"}
	list := func() (groupMemberRoles, error) {
		return groupMemberRoles{"a@example.com": "OWNER"}, nil
	}

	err := waitForGroupMembers(context.Background(), 2*time.Second, expected, list)
	if err == nil || !IsNotConsistent(err) {
		t.Fatalf("expected a consistency timeout, got %v", err)
	}
}

func TestWaitForGroupMembers_listErrorIsNotRetried(t *testing.T) {
	t.Parallel()

	calls := 0
	list := func() (groupMemberRoles, error) {
		calls++
		return nil, errors.New("boom")
	}

	err := waitForGroupMembers(context.Background(), 20*time.Second, groupMemberRoles{}, list)
	if err == nil || IsNotConsistent(err) {
		t.Fatalf("expected the listing error to surface, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected one listing, got %d", calls)
	}
}

func TestGroupMembersMatch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		expected, listed groupMemberRoles
		want             bool
	}{
		{"equal", groupMemberRoles{"a@x": "OWNER"}, groupMemberRoles{"a@x": "OWNER"}, true},
		{"missing member", groupMemberRoles{"a@x": "OWNER", "b@x": "MEMBER"}, groupMemberRoles{"a@x": "OWNER"}, false},
		{"extra member", groupMemberRoles{"a@x": "OWNER"}, groupMemberRoles{"a@x": "OWNER", "b@x": "MEMBER"}, false},
		{"stale role", groupMemberRoles{"a@x": "MEMBER"}, groupMemberRoles{"a@x": "OWNER"}, false},
		{"both empty", groupMemberRoles{}, groupMemberRoles{}, true},
	}
	for _, c := range cases {
		if got := groupMembersMatch(c.expected, c.listed); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
