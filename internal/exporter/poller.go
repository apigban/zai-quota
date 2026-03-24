package exporter

import (
	"context"
	"log"
	"time"

	"zai-quota/internal/models"
)

const MinimumPollInterval = 60 * time.Second

type QuotaFetcher interface {
	FetchQuota(ctx context.Context) (*models.QuotaResponse, error)
}

type Poller struct {
	client      QuotaFetcher
	cache       *CachedMetrics
	interval    time.Duration
	lastPoll    time.Time
	pollingDone chan struct{}
}

func NewPoller(client QuotaFetcher, cache *CachedMetrics, interval time.Duration) *Poller {
	if interval < MinimumPollInterval {
		interval = MinimumPollInterval
	}

	return &Poller{
		client:      client,
		cache:       cache,
		interval:    interval,
		pollingDone: make(chan struct{}),
	}
}

func (p *Poller) Start(ctx context.Context) {
	go p.run(ctx)
}

func (p *Poller) Stop() {
	<-p.pollingDone
}

func (p *Poller) run(ctx context.Context) {
	defer close(p.pollingDone)

	p.poll(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	start := time.Now()
	quota, err := p.client.FetchQuota(ctx)
	duration := time.Since(start)

	if err != nil {
		log.Printf("exporter: failed to fetch quota: %v", err)
		p.cache.Update(nil, duration, false)
		return
	}

	p.cache.Update(quota, duration, true)
	p.lastPoll = time.Now()
}
