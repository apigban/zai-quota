package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"zai-quota/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		endpoint       string
		timeoutSeconds int
		wantTimeout    time.Duration
	}{
		{
			name:           "default timeout",
			apiKey:         "test-key",
			endpoint:       "https://api.example.com",
			timeoutSeconds: 0, // should default to 5
			wantTimeout:    5 * time.Second,
		},
		{
			name:           "custom timeout",
			apiKey:         "test-key",
			endpoint:       "https://api.example.com",
			timeoutSeconds: 10,
			wantTimeout:    10 * time.Second,
		},
		{
			name:           "negative timeout defaults to 5",
			apiKey:         "test-key",
			endpoint:       "https://api.example.com",
			timeoutSeconds: -1,
			wantTimeout:    5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.endpoint, tt.timeoutSeconds)

			assert.NotNil(t, client)
			assert.Equal(t, tt.apiKey, client.apiKey)
			assert.Equal(t, tt.endpoint, client.endpoint)
			assert.NotNil(t, client.httpClient)
			assert.Equal(t, tt.wantTimeout, client.httpClient.Timeout)

			// Verify transport configuration
			transport, ok := client.httpClient.Transport.(*http.Transport)
			require.True(t, ok, "transport should be *http.Transport")
			assert.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout)
			assert.Equal(t, 10*time.Second, transport.ResponseHeaderTimeout)
		})
	}
}

func TestFetchQuota_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return success response with proper structure
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"code": 200,
			"msg": "OK",
			"data": {
				"limits": [
					{
						"type": "TIME_LIMIT",
						"percentage": 50,
						"usage": 1000,
						"currentValue": 500,
						"remaining": 500,
						"nextResetTime": 1709856000
					}
				],
				"level": "pro"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Limits, 1)
	assert.Equal(t, "TIME_LIMIT", resp.Limits[0].Type)
	assert.Equal(t, 1000, resp.Limits[0].Usage)
	assert.Equal(t, 500, resp.Limits[0].CurrentValue)
	assert.Equal(t, 500, resp.Limits[0].Remaining)
	assert.Equal(t, "pro", resp.Level)
}

func TestFetchQuota_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid API key"}`))
	}))
	defer server.Close()

	client := NewClient("invalid-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitAuth, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "unauthorized"))
}

func TestFetchQuota_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "access denied"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitAuth, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "forbidden"))
}

func TestFetchQuota_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitError, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "server error"))
}

func TestFetchQuota_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay response longer than client timeout
		time.Sleep(6 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"limit": 1000, "remaining": 500}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, 1) // 1 second timeout
	ctx := context.Background()

	start := time.Now()
	resp, err := client.FetchQuota(ctx)
	duration := time.Since(start)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitNetwork, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "timeout") || strings.Contains(qErr.Message, "failed"))

	// Verify timeout was enforced (should take ~1s, not 6s)
	assert.True(t, duration < 3*time.Second, "request should timeout quickly")
}

func TestFetchQuota_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slow response that won't complete
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, 10)

	// Create a context that will be canceled quickly
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitNetwork, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "canceled") || strings.Contains(qErr.Message, "failed"))
}

func TestFetchQuota_NetworkError(t *testing.T) {
	// Use an invalid URL to trigger a network error
	client := NewClient("test-key", "http://localhost:99999", 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitNetwork, qErr.Code)
}

func TestFetchQuota_UnexpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // 404
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)

	qErr, ok := err.(*errors.QuotaError)
	require.True(t, ok, "error should be *errors.QuotaError")
	assert.Equal(t, errors.ExitError, qErr.Code)
	assert.True(t, strings.Contains(qErr.Message, "unexpected"))
}

func TestFetchQuota_NonSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": false, "code": 400, "msg": "Invalid request parameters", "data": null}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	assert.Contains(t, err.Error(), "Invalid request parameters")
}

func TestFetchQuota_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, invalid json`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse response")
}

func TestFetchQuota_MissingRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`)) // Empty response - missing all required fields
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid response structure")
	assert.Contains(t, err.Error(), "missing required fields")
}

func TestFetchQuota_MissingDataField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "code": 200, "msg": "OK", "data": null}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	assert.Nil(t, resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing data field")
}

func TestFetchQuota_ParsingWithEmptyReset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"code": 200,
			"msg": "OK",
			"data": {
				"limits": [
					{
						"type": "TIME_LIMIT",
						"percentage": 50,
						"usage": 1000,
						"currentValue": 500,
						"remaining": 500,
						"nextResetTime": 0
					}
				],
				"level": "pro"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key", server.URL, 5)
	ctx := context.Background()

	resp, err := client.FetchQuota(ctx)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Limits, 1)
	assert.Equal(t, "TIME_LIMIT", resp.Limits[0].Type)
	assert.Equal(t, 1000, resp.Limits[0].Usage)
	assert.Equal(t, 500, resp.Limits[0].CurrentValue)
	assert.Equal(t, 500, resp.Limits[0].Remaining)
	assert.Equal(t, int64(0), resp.Limits[0].NextResetTime)
	assert.Equal(t, "pro", resp.Level)
}
