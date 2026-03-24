package exporter

import (
	"context"
	"testing"
	"time"

	"zai-quota/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewPoller_EnforcesMinimumInterval(t *testing.T) {
	cache := NewCachedMetrics()
	client := newMockClient()

	poller := NewPoller(client, cache, 30*time.Second)
	assert.Equal(t, MinimumPollInterval, poller.interval)
}

func TestNewPoller_RespectsValidInterval(t *testing.T) {
	cache := NewCachedMetrics()
	client := newMockClient()

	poller := NewPoller(client, cache, 120*time.Second)
	assert.Equal(t, 120*time.Second, poller.interval)
}

func TestPoller_Poll(t *testing.T) {
	cache := NewCachedMetrics()
	client := newMockClient()
	poller := NewPoller(client, cache, 60*time.Second)

	poller.poll(context.Background())

	limits, level, _, _, success := cache.Get()
	assert.True(t, success)
	assert.Equal(t, "pro", level)
	assert.Len(t, limits, 2)
}

func TestPoller_PollWithError(t *testing.T) {
	cache := NewCachedMetrics()
	client := newMockClientWithError()
	poller := NewPoller(client, cache, 60*time.Second)

	poller.poll(context.Background())

	_, _, _, _, success := cache.Get()
	assert.False(t, success)
}

type mockClient struct {
	quota *models.QuotaResponse
	err   error
}

func newMockClient() *mockClient {
	return &mockClient{
		quota: &models.QuotaResponse{
			Limits: []models.Limit{
				{Type: "TOKENS_LIMIT", Percentage: 50, NextResetTime: time.Now().Add(5 * time.Hour).UnixMilli()},
				{Type: "TIME_LIMIT", CurrentValue: 10, Usage: 1000, Remaining: 990, NextResetTime: time.Now().Add(24 * time.Hour).UnixMilli()},
			},
			Level: "pro",
		},
	}
}

func newMockClientWithError() *mockClient {
	return &mockClient{
		err: assert.AnError,
	}
}

func (m *mockClient) FetchQuota(ctx context.Context) (*models.QuotaResponse, error) {
	return m.quota, m.err
}
