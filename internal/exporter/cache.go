package exporter

import (
	"sync"
	"time"

	"zai-quota/internal/models"
)

type CachedMetrics struct {
	Limits            []models.Limit
	Level             string
	LastScrapeTime    time.Time
	ScrapeDuration    time.Duration
	LastScrapeSuccess bool
	mu                sync.RWMutex
}

func NewCachedMetrics() *CachedMetrics {
	return &CachedMetrics{
		LastScrapeSuccess: false,
	}
}

func (c *CachedMetrics) Update(quota *models.QuotaResponse, duration time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if quota != nil {
		c.Limits = quota.Limits
		c.Level = quota.Level
	}
	c.LastScrapeTime = time.Now()
	c.ScrapeDuration = duration
	c.LastScrapeSuccess = success
}

func (c *CachedMetrics) Get() ([]models.Limit, string, time.Time, time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Limits, c.Level, c.LastScrapeTime, c.ScrapeDuration, c.LastScrapeSuccess
}

func (c *CachedMetrics) HasData() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.Limits) > 0 || c.LastScrapeTime.After(time.Time{})
}

func (c *CachedMetrics) GetPromptLimit() *models.Limit {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.Limits {
		if c.Limits[i].Type == "TOKENS_LIMIT" {
			return &c.Limits[i]
		}
	}
	return nil
}

func (c *CachedMetrics) GetTimeLimit() *models.Limit {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.Limits {
		if c.Limits[i].Type == "TIME_LIMIT" {
			return &c.Limits[i]
		}
	}
	return nil
}
