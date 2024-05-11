// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
)

type RetryErrorPredicateFunc func(error) (bool, string)

/** ADD GLOBAL ERROR RETRY PREDICATES HERE **/
// Retry predicates that shoud apply to all requests should be added here.
var defaultErrorRetryPredicates = []RetryErrorPredicateFunc{
	// Common network errors (usually wrapped by URL error)
	isNetworkTemporaryError,
	isNetworkTimeoutError,
	isIoEOFError,
	isConnectionResetNetworkError,

	// Common error codes
	isCommonRetryableErrorCode,
	isRateLimitExceeded,
}

/** END GLOBAL ERROR RETRY PREDICATES HERE **/

func isNetworkTemporaryError(err error) (bool, string) {
	if netErr, ok := err.(*net.OpError); ok && netErr.Temporary() {
		return true, "marked as timeout"
	}
	if urlerr, ok := err.(*url.Error); ok && urlerr.Temporary() {
		return true, "marked as timeout"
	}
	return false, ""
}

func isNetworkTimeoutError(err error) (bool, string) {
	if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
		return true, "marked as timeout"
	}
	if urlerr, ok := err.(*url.Error); ok && urlerr.Timeout() {
		return true, "marked as timeout"
	}
	return false, ""
}

func isIoEOFError(err error) (bool, string) {
	if err == io.ErrUnexpectedEOF {
		return true, "Got unexpected EOF"
	}

	if urlerr, urlok := err.(*url.Error); urlok {
		wrappedErr := urlerr.Unwrap()
		if wrappedErr == io.ErrUnexpectedEOF {
			return true, "Got unexpected EOF"
		}
	}
	return false, ""
}

const connectionResetByPeerErr = ": connection reset by peer"

func isConnectionResetNetworkError(err error) (bool, string) {
	if strings.HasSuffix(err.Error(), connectionResetByPeerErr) {
		return true, fmt.Sprintf("reset connection error: %v", err)
	}
	return false, ""
}

// Retry on common googleapi error codes for retryable errors.
// what retryable error codes apply to which API.
func isCommonRetryableErrorCode(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503 {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}

	if gerr.Code == 401 && strings.Contains(gerr.Body, "Login Required") {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}

	// Unfortunately, the Google API sometimes returns 403 - Not Authorized to access this resource/api after a create operation
	// even though the resource was created successfully. Becasue of this, we should retry on 403 errors as well
	// This will lead to slower error responses when there is an actual problem, but it is better than failing the operation
	// when the resource was created successfully. Some of the actual problems include:
	// - trying to modify certain fields on an admin user without domain-wide delegation
	// - trying to undertake operations without the correct permissions
	if gerr.Code == 403 && strings.Contains(gerr.Body, "Not Authorized to access this resource/api") {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}
	return false, ""
}

func isRateLimitExceeded(err error) (bool, string) {
	gerr, ok := err.(*googleapi.Error)
	if !ok {
		return false, ""
	}

	if gerr.Code == 429 {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s, delaying retry by 10 seconds", err)
		time.Sleep(10 * time.Second)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}

	if gerr.Code == 403 && (strings.Contains(gerr.Error(), "Quota exceeded") || strings.Contains(gerr.Error(), "quotaExceeded")) {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s, , delaying retry by 10 seconds", err)
		time.Sleep(10 * time.Second)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}

	if gerr.Code == 403 {
		log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
		return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
	}

	return false, ""
}

// IsNotFound reports whether err is the result of the
// server replying with http.StatusNotFound.
// Such error values are sometimes returned by "Do" methods
// on calls when creation of ressource was too recent to return values
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*googleapi.Error)
	return ok && ae.Code == http.StatusNotFound
}
