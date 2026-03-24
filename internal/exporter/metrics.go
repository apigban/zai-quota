package exporter

import (
	"time"

	"zai-quota/internal/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "zai_quota"

type Metrics struct {
	PromptUsageRatio        prometheus.Gauge
	PromptResetTimestamp    prometheus.Gauge
	ToolCallsUsed           prometheus.Gauge
	ToolCallsLimit          prometheus.Gauge
	ToolCallsRemaining      prometheus.Gauge
	ToolCallsResetTimestamp prometheus.Gauge
	ToolCallsByTool         *prometheus.GaugeVec
	Info                    *prometheus.GaugeVec
	Up                      prometheus.Gauge
	LastScrapeTimestamp     prometheus.Gauge
	ScrapeDuration          prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		PromptUsageRatio: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "prompt_usage_ratio",
			Help:      "Current prompt usage as a ratio (0-1) of the 5-hour limit",
		}),
		PromptResetTimestamp: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "prompt_reset_timestamp_seconds",
			Help:      "Unix timestamp when the prompt limit resets",
		}),
		ToolCallsUsed: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tool_calls_used",
			Help:      "Number of tool calls used in the current period",
		}),
		ToolCallsLimit: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tool_calls_limit",
			Help:      "Maximum number of tool calls allowed in the period",
		}),
		ToolCallsRemaining: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tool_calls_remaining",
			Help:      "Number of tool calls remaining in the current period",
		}),
		ToolCallsResetTimestamp: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tool_calls_reset_timestamp_seconds",
			Help:      "Unix timestamp when the tool call limit resets",
		}),
		ToolCallsByTool: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tool_calls_by_tool",
			Help:      "Number of tool calls per tool",
		}, []string{"tool"}),
		Info: promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "info",
			Help:      "Information about the Z.ai subscription",
		}, []string{"level"}),
		Up: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the last scrape was successful (1) or failed (0)",
		}),
		LastScrapeTimestamp: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "last_scrape_timestamp_seconds",
			Help:      "Unix timestamp of the last scrape from Z.ai API",
		}),
		ScrapeDuration: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "scrape_duration_seconds",
			Help:      "Duration of the last scrape from Z.ai API in seconds",
		}),
	}

	return m
}

func (m *Metrics) Update(limits []models.Limit, level string, lastScrape time.Time, duration time.Duration, success bool) {
	if success {
		m.Up.Set(1)
	} else {
		m.Up.Set(0)
	}

	m.LastScrapeTimestamp.Set(float64(lastScrape.Unix()))
	m.ScrapeDuration.Set(duration.Seconds())

	m.Info.Reset()
	m.Info.WithLabelValues(level).Set(1)

	for i := range limits {
		limit := &limits[i]
		switch limit.Type {
		case "TOKENS_LIMIT":
			m.PromptUsageRatio.Set(float64(limit.Percentage) / 100.0)
			m.PromptResetTimestamp.Set(float64(limit.NextResetTime) / 1000.0)

		case "TIME_LIMIT":
			m.ToolCallsUsed.Set(float64(limit.CurrentValue))
			m.ToolCallsLimit.Set(float64(limit.Usage))
			m.ToolCallsRemaining.Set(float64(limit.Remaining))
			m.ToolCallsResetTimestamp.Set(float64(limit.NextResetTime) / 1000.0)

			m.ToolCallsByTool.Reset()
			for _, detail := range limit.UsageDetails {
				m.ToolCallsByTool.WithLabelValues(detail.ModelCode).Set(float64(detail.Usage))
			}
		}
	}
}
