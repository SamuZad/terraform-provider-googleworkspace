// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"errors"
	"strconv"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsCommonRetryableErrorCode_retryableErrorCode(t *testing.T) {
	codes := []int{500, 502, 503}
	for _, code := range codes {
		code := code
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			err := googleapi.Error{
				Code: code,
				Body: "some text describing error",
			}
			isRetryable, _ := isCommonRetryableErrorCode(&err)
			if !isRetryable {
				t.Errorf("Error not detected as retryable")
			}
		})
	}
}

func TestIsCommonRetryableErrorCode_otherError(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Some unretryable issue",
	}
	isRetryable, _ := isCommonRetryableErrorCode(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsOperationReadQuotaError_quotaExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Request rate higher than configured., quotaExceeded",
	}
	isRetryable, _ := isRateLimitExceeded(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsOperationReadQuotaError_rateLimitExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 429,
		Body: "Rate Limit Exceeded., rateLimitExceeded",
	}
	isRetryable, _ := isRateLimitExceeded(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsConcurrentUpdateError_retryable(t *testing.T) {
	cases := map[string]googleapi.Error{
		"412 precondition failed": {
			Code: 412,
			Body: "Precondition Failed",
		},
		"400 invalid input resource_id": {
			Code: 400,
			Body: "Invalid Input: resource_id",
		},
	}
	for name, err := range cases {
		err := err
		t.Run(name, func(t *testing.T) {
			isRetryable, _ := isConcurrentUpdateError(&err)
			if !isRetryable {
				t.Errorf("Error not detected as retryable")
			}
		})
	}
}

func TestIsConcurrentUpdateError_notRetryable(t *testing.T) {
	cases := map[string]googleapi.Error{
		"400 unrelated body": {
			Code: 400,
			Body: "Invalid Input: some other field",
		},
		"404 not found": {
			Code: 404,
			Body: "Not Found",
		},
	}
	for name, err := range cases {
		err := err
		t.Run(name, func(t *testing.T) {
			isRetryable, _ := isConcurrentUpdateError(&err)
			if isRetryable {
				t.Errorf("Error incorrectly detected as retryable")
			}
		})
	}
}

func TestIsConcurrentUpdateError_nonGoogleError(t *testing.T) {
	isRetryable, _ := isConcurrentUpdateError(errors.New("some non-google error"))
	if isRetryable {
		t.Errorf("Non-Google error incorrectly detected as retryable")
	}
}

func TestGoogle404Error(t *testing.T) {
	gerr := googleapi.Error{
		Code:    404,
		Message: "notfound",
	}
	err := &gerr

	expected := true

	if isNotFound(err) != expected {
		t.Error("Failed: The error should have been detected as 404")
	}
}

func TestGoogleNot404Error(t *testing.T) {
	gerr := googleapi.Error{
		Code:    200,
		Message: "notfound",
	}
	err := &gerr

	expected := false

	if isNotFound(err) != expected {
		t.Error("Failed: The error was detected as a 404 but should not have been")
	}
}
