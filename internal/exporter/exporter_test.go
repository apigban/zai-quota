package exporter

import (
	"context"
	"sync"
	"testing"
	"time"

	"zai-quota/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuotaFetcher struct {
	mu        sync.Mutex
	quota     *models.QuotaResponse
	err       error
	callCount int
}

func (m *mockQuotaFetcher) FetchQuota(ctx context.Context) (*models.QuotaResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	return m.quota, m.err
}

func (m *mockQuotaFetcher) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestNewExporter_Validation(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *ExporterConfig
		expectedErr string
	}{
		{
			name: "poll interval too short",
			cfg: &ExporterConfig{
				APIKey:         "test-key",
				Endpoint:       "https://api.example.com",
				TimeoutSeconds: 30,
				PollInterval:   30 * time.Second,
				ListenAddr:     ":9090",
			},
			expectedErr: "poll interval must be at least",
		},
		{
			name: "valid config",
			cfg: &ExporterConfig{
				APIKey:         "test-key",
				Endpoint:       "https://api.example.com",
				TimeoutSeconds: 30,
				PollInterval:   60 * time.Second,
				ListenAddr:     ":9090",
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exporter, err := NewExporter(tt.cfg)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, exporter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, exporter)
			}
		})
	}
}

func TestExporter_Lifecycle(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		quota: &models.QuotaResponse{
			Level: "pro",
			Limits: []models.Limit{
				{Type: "TOKENS_LIMIT", Percentage: 50, NextResetTime: time.Now().Add(5 * time.Hour).UnixMilli()},
				{Type: "TIME_LIMIT", CurrentValue: 10, Usage: 1000, Remaining: 990, NextResetTime: time.Now().Add(24 * time.Hour).UnixMilli()},
			},
		},
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)
	require.NotNil(t, exporter)

	exporter.client = mockClient

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(500 * time.Millisecond)

	assert.True(t, exporter.cache.HasData(), "Cache should have data after running")

	err = <-errChan
	assert.NoError(t, err, "Exporter should shut down cleanly")
}

func TestExporter_InitialPoll(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		quota: &models.QuotaResponse{
			Level: "pro",
			Limits: []models.Limit{
				{Type: "TOKENS_LIMIT", Percentage: 50},
			},
		},
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)

	exporter.client = mockClient

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	assert.GreaterOrEqual(t, mockClient.GetCallCount(), 1, "Should have polled at least once")

	<-errChan
}

func TestExporter_HandlesAPIErrors(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		err: assert.AnError,
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)

	exporter.client = mockClient

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	limits, level, _, _, success := exporter.cache.Get()
	assert.False(t, success, "Cache should indicate failed scrape")
	assert.Empty(t, limits)
	assert.Empty(t, level)

	<-errChan
}

func TestExporter_GracefulShutdown(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		quota: &models.QuotaResponse{
			Level:  "pro",
			Limits: []models.Limit{},
		},
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)

	exporter.client = mockClient

	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case err := <-errChan:
		assert.NoError(t, err, "Should shut down without error")
	case <-time.After(2 * time.Second):
		t.Fatal("Exporter did not shut down within timeout")
	}
}

func TestExporter_MetricsUpdated(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		quota: &models.QuotaResponse{
			Level: "pro",
			Limits: []models.Limit{
				{Type: "TOKENS_LIMIT", Percentage: 75},
			},
		},
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)

	exporter.client = mockClient

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	limits, _, _, _, _ := exporter.cache.Get()
	assert.NotEmpty(t, limits, "Metrics should be updated with data from API")

	<-errChan
}

func TestExporter_HTTPServerRunning(t *testing.T) {
	mockClient := &mockQuotaFetcher{
		quota: &models.QuotaResponse{
			Level:  "pro",
			Limits: []models.Limit{},
		},
	}

	cfg := &ExporterConfig{
		APIKey:         "test-key",
		Endpoint:       "https://api.example.com",
		TimeoutSeconds: 30,
		PollInterval:   60 * time.Second,
		ListenAddr:     ":0",
	}

	exporter, err := NewExporter(cfg)
	require.NoError(t, err)

	exporter.client = mockClient

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- exporter.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	assert.NotNil(t, exporter.httpServer, "HTTP server should be initialized")

	<-errChan
}
