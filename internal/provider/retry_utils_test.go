// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/api/googleapi"
)

func stubNotFoundBackoff(t *testing.T, backoff []time.Duration) {
	t.Helper()
	orig := notFoundBackoff
	notFoundBackoff = backoff
	t.Cleanup(func() { notFoundBackoff = orig })
}

func TestRetryNotFound_succeedsFirstTry(t *testing.T) {
	stubNotFoundBackoff(t, []time.Duration{0, 0, 0, 0, 0})

	calls := 0
	err := retryNotFound(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryNotFound_recoversAfter404s(t *testing.T) {
	stubNotFoundBackoff(t, []time.Duration{0, 0, 0, 0, 0})

	calls := 0
	err := retryNotFound(context.Background(), func() error {
		calls++
		if calls <= 2 {
			return &googleapi.Error{Code: 404, Body: "not found"}
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryNotFound_exhaustsRetriesOnPersistent404(t *testing.T) {
	stubNotFoundBackoff(t, []time.Duration{0, 0, 0, 0, 0})

	calls := 0
	err := retryNotFound(context.Background(), func() error {
		calls++
		return &googleapi.Error{Code: 404, Body: "not found"}
	})
	if !isApiErrorWithCode(err, 404) {
		t.Errorf("expected 404 error after exhausting retries, got: %v", err)
	}
	if calls != 6 {
		t.Errorf("expected 6 calls (1 + 5 retries), got %d", calls)
	}
}

func TestRetryNotFound_non404NotRetried(t *testing.T) {
	stubNotFoundBackoff(t, []time.Duration{0, 0, 0, 0, 0})

	calls := 0
	err := retryNotFound(context.Background(), func() error {
		calls++
		return &googleapi.Error{Code: 403, Body: "forbidden"}
	})
	if !isApiErrorWithCode(err, 403) {
		t.Errorf("expected 403 error, got: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryNotFound_contextCancelled(t *testing.T) {
	stubNotFoundBackoff(t, []time.Duration{10})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retryNotFound(ctx, func() error {
		calls++
		return &googleapi.Error{Code: 404, Body: "not found"}
	})
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestIsNotConsistent_retryable(t *testing.T) {
	err := fmt.Errorf("something timed out while waiting")
	isRetryable := IsNotConsistent(err)
	if !isRetryable {
		t.Errorf("inconsistent error not detected as temporarily unavailable")
	}
}

func TestIsNotConsistent_other(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "some text describing error",
	}
	isRetryable := IsNotConsistent(&err)
	if isRetryable {
		t.Errorf("404 error detected as inconsistency error")
	}
}
