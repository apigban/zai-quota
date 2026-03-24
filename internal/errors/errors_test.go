package errors

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

)

func TestExitCodeConstants(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"ExitSuccess", ExitSuccess, 0},
		{"ExitError", ExitError, 1},
		{"ExitNetwork", ExitNetwork, 2},
		{"ExitAuth", ExitAuth, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.want)
			}
		})
	}
}

func TestQuotaError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *QuotaError
		wantMsg string
	}{
		{
			name: "Error with underlying error",
			err: &QuotaError{
				Code:    ExitNetwork,
				Message: "network failed",
				Err:     errors.New("connection refused"),
			},
			wantMsg: "network failed: connection refused",
		},
		{
			name: "Error without underlying error",
			err: &QuotaError{
				Code:    ExitAuth,
				Message: "unauthorized",
				Err:     nil,
			},
			wantMsg: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestQuotaError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &QuotaError{
		Code:    ExitError,
		Message: "wrapped error",
		Err:     underlying,
	}

	if err.Unwrap() != underlying {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), underlying)
	}

	// Test with nil underlying error
	errNoUnwrap := &QuotaError{
		Code:    ExitError,
		Message: "no unwrap",
		Err:     nil,
	}

	if errNoUnwrap.Unwrap() != nil {
		t.Errorf("Unwrap() with nil Err = %v, want nil", errNoUnwrap.Unwrap())
	}
}

func TestQuotaError_ExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  *QuotaError
		want int
	}{
		{"Success exit code", &QuotaError{Code: ExitSuccess}, 0},
		{"Error exit code", &QuotaError{Code: ExitError}, 1},
		{"Network exit code", &QuotaError{Code: ExitNetwork}, 2},
		{"Auth exit code", &QuotaError{Code: ExitAuth}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.ExitCode(); got != tt.want {
				t.Errorf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNewNetworkError(t *testing.T) {
	underlying := errors.New("timeout")
	err := NewNetworkError("network timeout", underlying)

	if err.Code != ExitNetwork {
		t.Errorf("NewNetworkError() Code = %d, want %d", err.Code, ExitNetwork)
	}
	if err.Message != "network timeout" {
		t.Errorf("NewNetworkError() Message = %q, want %q", err.Message, "network timeout")
	}
	if err.Err != underlying {
		t.Errorf("NewNetworkError() Err = %v, want %v", err.Err, underlying)
	}

	// Verify Error() method works
	if !strings.Contains(err.Error(), "network timeout") {
		t.Errorf("NewNetworkError() Error() should contain message")
	}
}

func TestNewAuthError(t *testing.T) {
	underlying := errors.New("invalid token")
	err := NewAuthError("authentication failed", underlying)

	if err.Code != ExitAuth {
		t.Errorf("NewAuthError() Code = %d, want %d", err.Code, ExitAuth)
	}
	if err.Message != "authentication failed" {
		t.Errorf("NewAuthError() Message = %q, want %q", err.Message, "authentication failed")
	}
	if err.Err != underlying {
		t.Errorf("NewAuthError() Err = %v, want %v", err.Err, underlying)
	}

	// Verify Error() method works
	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("NewAuthError() Error() should contain message")
	}
}

func TestNewGenericError(t *testing.T) {
	underlying := errors.New("something went wrong")
	err := NewGenericError("generic error", underlying)

	if err.Code != ExitError {
		t.Errorf("NewGenericError() Code = %d, want %d", err.Code, ExitError)
	}
	if err.Message != "generic error" {
		t.Errorf("NewGenericError() Message = %q, want %q", err.Message, "generic error")
	}
	if err.Err != underlying {
		t.Errorf("NewGenericError() Err = %v, want %v", err.Err, underlying)
	}
}

func TestMapHTTPStatusToExitCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       int
	}{
		{"Unauthorized maps to ExitAuth", http.StatusUnauthorized, ExitAuth},
		{"Forbidden maps to ExitAuth", http.StatusForbidden, ExitAuth},
		{"NotFound maps to ExitError", http.StatusNotFound, ExitError},
		{"InternalServerError maps to ExitError", http.StatusInternalServerError, ExitError},
		{"BadRequest maps to ExitError", http.StatusBadRequest, ExitError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapHTTPStatusToExitCode(tt.statusCode)
			if got != tt.want {
				t.Errorf("MapHTTPStatusToExitCode(%d) = %d, want %d", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "Nil error returns ExitSuccess",
			err:  nil,
			want: ExitSuccess,
		},
		{
			name: "Timeout error returns ExitNetwork",
			err:  &net.OpError{Err: &timeoutError{}},
			want: ExitNetwork,
		},
		{
			name: "DNS error returns ExitNetwork",
			err:  &net.DNSError{Err: "no such host"},
			want: ExitNetwork,
		},
		{
			name: "Network error returns ExitNetwork",
			err:  &net.OpError{Op: "dial"},
			want: ExitNetwork,
		},

		{
			name: "QuotaError with ExitNetwork code",
			err:  NewNetworkError("test", nil),
			want: ExitNetwork,
		},
		{
			name: "QuotaError with ExitAuth code",
			err:  NewAuthError("test", nil),
			want: ExitAuth,
		},
		{
			name: "QuotaError with ExitError code",
			err:  NewGenericError("test", nil),
			want: ExitError,
		},
		{
			name: "Generic error returns ExitError",
			err:  errors.New("some error"),
			want: ExitError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyError(tt.err)
			if got != tt.want {
				t.Errorf("ClassifyError() = %d, want %d", got, tt.want)
			}
		})
	}
}

// Helper type for timeout error testing
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

func TestErrorInterface(t *testing.T) {
	var _ error = NewNetworkError("test", nil)
	var _ error = NewAuthError("test", nil)
	var _ error = NewGenericError("test", nil)

	err := NewNetworkError("test", nil)
	if err.Error() == "" {
		t.Error("Error() should return non-empty string")
	}
}

func TestErrorWrapping(t *testing.T) {
	underlying1 := errors.New("level 1")
	underlying2 := fmt.Errorf("level 2: %w", underlying1)
	wrapped := NewNetworkError("network error", underlying2)

	// Test that all errors can be unwrapped
	err := errors.Unwrap(wrapped)
	if err == nil {
		t.Error("Unwrap should return an error")
	}

	// Test errors.Is for underlying error
	if !errors.Is(wrapped, underlying2) {
		t.Error("errors.Is should find underlying error")
	}
}
