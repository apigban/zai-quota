package exporter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestLandingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	landingHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Z.ai Quota Exporter")
	assert.Contains(t, w.Body.String(), "/metrics")
	assert.Contains(t, w.Body.String(), "/health")
}

func TestNewServerHandler_Routes(t *testing.T) {
	reg := prometheus.NewRegistry()
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	handler := NewServerHandler(metricsHandler)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health endpoint",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "landing page",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Z.ai Quota Exporter",
		},
		{
			name:           "metrics endpoint",
			path:           "/metrics",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestNewServerHandler_MetricsFormat(t *testing.T) {
	reg := prometheus.NewRegistry()
	_ = NewMetrics(reg)
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	handler := NewServerHandler(metricsHandler)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")
	assert.Contains(t, w.Body.String(), "# HELP")
	assert.Contains(t, w.Body.String(), "# TYPE")
}

func TestNewServerHandler_MethodNotAllowed(t *testing.T) {
	reg := prometheus.NewRegistry()
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	handler := NewServerHandler(metricsHandler)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/metrics", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)
		})
	}
}

func TestLandingPage_ContainsMetricNames(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	landingHandler(w, req)

	body := w.Body.String()

	expectedMetrics := []string{
		"zai_quota_prompt_usage_ratio",
		"zai_quota_tool_calls_used",
		"zai_quota_tool_calls_limit",
		"zai_quota_tool_calls_by_tool",
		"zai_quota_up",
	}

	for _, metric := range expectedMetrics {
		assert.Contains(t, body, metric, "Landing page should mention metric: %s", metric)
	}
}

func TestLandingPage_HTMLStructure(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	landingHandler(w, req)

	body := w.Body.String()

	assert.True(t, strings.Contains(body, "<!DOCTYPE html>"), "Should have DOCTYPE")
	assert.True(t, strings.Contains(body, "<html>"), "Should have html tag")
	assert.True(t, strings.Contains(body, "</html>"), "Should close html tag")
	assert.True(t, strings.Contains(body, "<head>"), "Should have head tag")
	assert.True(t, strings.Contains(body, "<body>"), "Should have body tag")
	assert.True(t, strings.Contains(body, "<h1>"), "Should have h1 tag")
	assert.True(t, strings.Contains(body, "<a href"), "Should have links")
}
