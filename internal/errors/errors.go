package errors

import (
	"fmt"
	"net"
	"net/http"

)

const (
	ExitSuccess = 0
	ExitError   = 1 // Generic error
	ExitNetwork = 2 // Network/timeout error
	ExitAuth    = 3 // Authentication error
)

// QuotaError represents an error with an associated exit code
type QuotaError struct {
	Code    int
	Message string
	Err     error
}

// NewNetworkError creates a new network error with ExitNetwork exit code
func NewNetworkError(msg string, err error) *QuotaError {
	return &QuotaError{
		Code:    ExitNetwork,
		Message: msg,
		Err:     err,
	}
}

// NewAuthError creates a new authentication error with ExitAuth exit code
func NewAuthError(msg string, err error) *QuotaError {
	return &QuotaError{
		Code:    ExitAuth,
		Message: msg,
		Err:     err,
	}
}

// NewGenericError creates a new generic error with ExitError exit code
func NewGenericError(msg string, err error) *QuotaError {
	return &QuotaError{
		Code:    ExitError,
		Message: msg,
		Err:     err,
	}
}

// Error implements the error interface
func (e *QuotaError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *QuotaError) Unwrap() error {
	return e.Err
}

// ExitCode returns the exit code for the error
func (e *QuotaError) ExitCode() int {
	return e.Code
}

// MapHTTPStatusToExitCode maps HTTP status codes to exit codes
func MapHTTPStatusToExitCode(statusCode int) int {
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return ExitAuth
	default:
		return ExitError
	}
}

// ClassifyError classifies an error and returns the appropriate exit code
func ClassifyError(err error) int {
	if err == nil {
		return ExitSuccess
	}

	// Check for network errors
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return ExitNetwork
	}

	// Check for other network-related errors
	if _, ok := err.(*net.DNSError); ok {
		return ExitNetwork
	}
	if _, ok := err.(*net.OpError); ok {
		return ExitNetwork
	}

	// Check for QuotaError
	if qErr, ok := err.(*QuotaError); ok {
		return qErr.Code
	}

	// Default to generic error
	return ExitError
}

