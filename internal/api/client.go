package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	apierrors "zai-quota/internal/errors"
	"zai-quota/internal/models"
)

// APIResponse represents the structure of the z.ai API response
type APIResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Msg     string                `json:"msg"`
	Data    *models.QuotaResponse `json:"data"`
}

// Client represents an HTTP client for the z.ai API
type Client struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
}

// NewClient creates a new API client with the specified configuration
func NewClient(apiKey, endpoint string, timeoutSeconds int) *Client {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 5 // default timeout
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   time.Duration(timeoutSeconds) * time.Second,
			Transport: transport,
		},
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

// Endpoint returns the API endpoint URL
func (c *Client) Endpoint() string {
	return c.endpoint
}

// FetchQuota retrieves quota information from the API
func (c *Client) FetchQuota(ctx context.Context) (*models.QuotaResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.endpoint, nil)
	if err != nil {
		return nil, apierrors.NewNetworkError("failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Check if it's a timeout error
		if ctx.Err() == context.Canceled {
			return nil, apierrors.NewNetworkError("request canceled", ctx.Err())
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, apierrors.NewNetworkError("request timeout", err)
		}
		return nil, apierrors.NewNetworkError("request failed", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, apierrors.NewAuthError("unauthorized: invalid API key", fmt.Errorf("status %d", resp.StatusCode))
		case http.StatusForbidden:
			return nil, apierrors.NewAuthError("forbidden: access denied", fmt.Errorf("status %d", resp.StatusCode))
		case http.StatusInternalServerError:
			return nil, apierrors.NewGenericError("server error", fmt.Errorf("status %d", resp.StatusCode))
		default:
			exitCode := apierrors.MapHTTPStatusToExitCode(resp.StatusCode)
			return nil, &apierrors.QuotaError{
				Code:    exitCode,
				Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			}
		}
	}

	// Parse and validate the API response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, apierrors.NewGenericError("failed to parse response", err)
	}

	// Validate response structure
	if apiResp.Code == 0 && apiResp.Msg == "" && apiResp.Data == nil {
		return nil, apierrors.NewGenericError("invalid response structure: missing required fields", fmt.Errorf("missing code, msg, or data"))
	}

	// Check for non-success response
	if !apiResp.Success {
		return nil, apierrors.NewGenericError(fmt.Sprintf("API error: %s", apiResp.Msg), fmt.Errorf("code: %d", apiResp.Code))
	}

	// Validate that the data field is present
	if apiResp.Data == nil {
		return nil, apierrors.NewGenericError("invalid response: missing data field", nil)
	}

	return apiResp.Data, nil
}
