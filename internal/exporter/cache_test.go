package exporter

import (
	"sync"
	"testing"
	"time"

	"zai-quota/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewCachedMetrics(t *testing.T) {
	cache := NewCachedMetrics()
	assert.NotNil(t, cache)
	assert.False(t, cache.LastScrapeSuccess)
	assert.Nil(t, cache.Limits)
	assert.Empty(t, cache.Level)
}

func TestCachedMetrics_Update(t *testing.T) {
	cache := NewCachedMetrics()
	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{Type: "TOKENS_LIMIT", Percentage: 50},
		},
		Level: "pro",
	}
	duration := 100 * time.Millisecond

	cache.Update(quota, duration, true)

	limits, level, _, dur, success := cache.Get()
	assert.Len(t, limits, 1)
	assert.Equal(t, "pro", level)
	assert.Equal(t, duration, dur)
	assert.True(t, success)
}

func TestCachedMetrics_Update_NilQuota(t *testing.T) {
	cache := NewCachedMetrics()
	duration := 50 * time.Millisecond

	cache.Update(nil, duration, false)

	_, _, _, dur, success := cache.Get()
	assert.Equal(t, duration, dur)
	assert.False(t, success)
}

func TestCachedMetrics_ConcurrentAccess(t *testing.T) {
	cache := NewCachedMetrics()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(val int) {
			defer wg.Done()
			quota := &models.QuotaResponse{
				Limits: []models.Limit{{Type: "TOKENS_LIMIT", Percentage: val}},
				Level:  "test",
			}
			cache.Update(quota, time.Millisecond, true)
		}(i)

		go func() {
			defer wg.Done()
			cache.Get()
		}()
	}

	wg.Wait()
}

func TestCachedMetrics_GetPromptLimit(t *testing.T) {
	cache := NewCachedMetrics()
	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{Type: "TIME_LIMIT", CurrentValue: 10},
			{Type: "TOKENS_LIMIT", Percentage: 75},
		},
		Level: "pro",
	}
	cache.Update(quota, time.Millisecond, true)

	limit := cache.GetPromptLimit()
	assert.NotNil(t, limit)
	assert.Equal(t, "TOKENS_LIMIT", limit.Type)
	assert.Equal(t, 75, limit.Percentage)
}

func TestCachedMetrics_GetTimeLimit(t *testing.T) {
	cache := NewCachedMetrics()
	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{Type: "TOKENS_LIMIT", Percentage: 75},
			{Type: "TIME_LIMIT", CurrentValue: 10},
		},
		Level: "pro",
	}
	cache.Update(quota, time.Millisecond, true)

	limit := cache.GetTimeLimit()
	assert.NotNil(t, limit)
	assert.Equal(t, "TIME_LIMIT", limit.Type)
	assert.Equal(t, 10, limit.CurrentValue)
}

func TestCachedMetrics_GetPromptLimit_NotFound(t *testing.T) {
	cache := NewCachedMetrics()
	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{Type: "TIME_LIMIT", CurrentValue: 10},
		},
		Level: "pro",
	}
	cache.Update(quota, time.Millisecond, true)

	limit := cache.GetPromptLimit()
	assert.Nil(t, limit)
}

func TestCachedMetrics_GetTimeLimit_NotFound(t *testing.T) {
	cache := NewCachedMetrics()
	quota := &models.QuotaResponse{
		Limits: []models.Limit{
			{Type: "TOKENS_LIMIT", Percentage: 75},
		},
		Level: "pro",
	}
	cache.Update(quota, time.Millisecond, true)

	limit := cache.GetTimeLimit()
	assert.Nil(t, limit)
}

func TestCachedMetrics_HasData(t *testing.T) {
	cache := NewCachedMetrics()
	assert.False(t, cache.HasData())

	quota := &models.QuotaResponse{
		Limits: []models.Limit{{Type: "TOKENS_LIMIT"}},
		Level:  "pro",
	}
	cache.Update(quota, time.Millisecond, true)
	assert.True(t, cache.HasData())
}
