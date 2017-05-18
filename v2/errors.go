package v2

import (
	"fmt"
	"net/http"
)

// HTTPStatusCodeError is an error type that provides additional information
// based on the Open Service Broker API conventions for returning information
// about errors.
//
// These errors may optionally provide a machine-readable error message and
// human-readable description.
//
// Checks for important errors in the API specification are implemented as
// utility methods:
//
// - IsConflictError
// - IsAsyncRequiredError
// - IsAppGUIDRequiredError
type HTTPStatusCodeError struct {
	StatusCode   int
	ErrorMessage *string
	Description  *string
}

func (e HTTPStatusCodeError) Error() string {
	var (
		errorMessage string = "nil"
		description  string = "nil"
	)
	if e.ErrorMessage != nil {
		errorMessage = *e.ErrorMessage
	}
	if e.Description != nil {
		description = *e.Description
	}

	return fmt.Sprintf("Status: %v; ErrorMessage: %v; Description: %v", e.StatusCode, errorMessage, description)
}

// IsConflictError returns whether the error represents a conflict.
func IsConflictError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	return statusCodeError.StatusCode == http.StatusConflict
}

const (
	asyncErrorMessage     = "AsyncRequired"
	asyncErrorDescription = "This service plan requires client support for asynchronous service operations."

	appGUIDRequiredErrorMessage     = "RequiresApp"
	appGUIDRequiredErrorDescription = "This service supports generation of credentials through binding an application only."
)

// IsAsyncRequiredError returns whether the error corresponds to the
// conventional way of indicating that a service requires asynchronous
// operations to perform an action.
func IsAsyncRequiredError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	if statusCodeError.StatusCode != http.StatusUnprocessableEntity {
		return false
	}

	if *statusCodeError.ErrorMessage != asyncErrorMessage {
		return false
	}

	return *statusCodeError.Description == asyncErrorDescription
}

// IsAppGUIDRequiredError returns whether the error corresponds to the
// conventional way of indicating that a service only supports credential-type
// bindings.
func IsAppGUIDRequiredError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	if statusCodeError.StatusCode != http.StatusUnprocessableEntity {
		return false
	}

	if *statusCodeError.ErrorMessage != appGUIDRequiredErrorMessage {
		return false
	}

	return *statusCodeError.Description == appGUIDRequiredErrorDescription
}
