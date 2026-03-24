package exporter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zai-quota/internal/api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ExporterConfig struct {
	APIKey         string
	Endpoint       string
	TimeoutSeconds int
	PollInterval   time.Duration
	ListenAddr     string
}

type Exporter struct {
	cfg        *ExporterConfig
	client     QuotaFetcher
	cache      *CachedMetrics
	metrics    *Metrics
	poller     *Poller
	registry   *prometheus.Registry
	httpServer *http.Server
}

func NewExporter(cfg *ExporterConfig) (*Exporter, error) {
	if cfg.PollInterval < MinimumPollInterval {
		return nil, fmt.Errorf("poll interval must be at least %v, got %v", MinimumPollInterval, cfg.PollInterval)
	}

	apiClient := api.NewClient(cfg.APIKey, cfg.Endpoint, cfg.TimeoutSeconds)
	cache := NewCachedMetrics()
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	return &Exporter{
		cfg:      cfg,
		client:   apiClient,
		cache:    cache,
		metrics:  metrics,
		registry: registry,
	}, nil
}

func (e *Exporter) Run(ctx context.Context) error {
	e.poller = NewPoller(e.client, e.cache, e.cfg.PollInterval)
	e.poller.Start(ctx)

	metricsHandler := promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{})
	handler := NewServerHandler(metricsHandler)

	e.httpServer = &http.Server{
		Addr:    e.cfg.ListenAddr,
		Handler: handler,
	}

	go e.updateMetricsLoop(ctx)

	errChan := make(chan error, 1)
	go func() {
		if err := e.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return e.httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

func (e *Exporter) updateMetricsLoop(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			limits, level, lastScrape, duration, success := e.cache.Get()
			e.metrics.Update(limits, level, lastScrape, duration, success)
		}
	}
}
