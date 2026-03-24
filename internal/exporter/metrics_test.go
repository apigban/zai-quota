package exporter

import (
	"strings"
	"testing"
	"time"

	"zai-quota/internal/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)
	assert.NotNil(t, metrics)
}

func TestMetrics_Update_Success(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type:          "TOKENS_LIMIT",
			Percentage:    28,
			NextResetTime: 1773234696431,
		},
		{
			Type:          "TIME_LIMIT",
			CurrentValue:  16,
			Usage:         1000,
			Remaining:     984,
			NextResetTime: 1775186469998,
			UsageDetails: []models.UsageDetail{
				{ModelCode: "search-prime", Usage: 4},
				{ModelCode: "web-reader", Usage: 3},
			},
		},
	}
	lastScrape := time.Date(2024, 3, 11, 12, 0, 0, 0, time.UTC)
	duration := 234 * time.Millisecond

	metrics.Update(limits, "pro", lastScrape, duration, true)

	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.Up))
	assert.Equal(t, float64(lastScrape.Unix()), testutil.ToFloat64(metrics.LastScrapeTimestamp))
	assert.Equal(t, duration.Seconds(), testutil.ToFloat64(metrics.ScrapeDuration))
	assert.Equal(t, 0.28, testutil.ToFloat64(metrics.PromptUsageRatio))
	assert.Equal(t, float64(1773234696.431), testutil.ToFloat64(metrics.PromptResetTimestamp))
	assert.Equal(t, float64(16), testutil.ToFloat64(metrics.ToolCallsUsed))
	assert.Equal(t, float64(1000), testutil.ToFloat64(metrics.ToolCallsLimit))
	assert.Equal(t, float64(984), testutil.ToFloat64(metrics.ToolCallsRemaining))
	assert.Equal(t, float64(1775186469.998), testutil.ToFloat64(metrics.ToolCallsResetTimestamp))
}

func TestMetrics_Update_Failure(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	lastScrape := time.Date(2024, 3, 11, 12, 0, 0, 0, time.UTC)
	duration := 100 * time.Millisecond

	metrics.Update(nil, "", lastScrape, duration, false)

	assert.Equal(t, float64(0), testutil.ToFloat64(metrics.Up))
}

func TestMetrics_Update_PercentageConversion(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type:       "TOKENS_LIMIT",
			Percentage: 50,
		},
	}

	metrics.Update(limits, "pro", time.Now(), time.Millisecond, true)

	assert.Equal(t, 0.5, testutil.ToFloat64(metrics.PromptUsageRatio))
}

func TestMetrics_Update_TimestampConversion(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type:          "TOKENS_LIMIT",
			NextResetTime: 1710123400000,
		},
	}

	metrics.Update(limits, "pro", time.Now(), time.Millisecond, true)

	expected := float64(1710123400)
	assert.Equal(t, expected, testutil.ToFloat64(metrics.PromptResetTimestamp))
}

func TestMetrics_Update_ToolCallsByTool(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type: "TIME_LIMIT",
			UsageDetails: []models.UsageDetail{
				{ModelCode: "search-prime", Usage: 4},
				{ModelCode: "web-reader", Usage: 3},
				{ModelCode: "doc-reader", Usage: 2},
			},
		},
	}

	metrics.Update(limits, "pro", time.Now(), time.Millisecond, true)

	metricFamilies, err := reg.Gather()
	assert.NoError(t, err)

	var found bool
	for _, mf := range metricFamilies {
		if mf.GetName() == "zai_quota_tool_calls_by_tool" {
			found = true
			assert.Len(t, mf.GetMetric(), 3)
		}
	}
	assert.True(t, found, "tool_calls_by_tool metric not found")
}

func TestMetrics_Update_EmptyUsageDetails(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type:         "TIME_LIMIT",
			UsageDetails: []models.UsageDetail{},
		},
	}

	metrics.Update(limits, "pro", time.Now(), time.Millisecond, true)

	metricFamilies, err := reg.Gather()
	assert.NoError(t, err)

	for _, mf := range metricFamilies {
		if mf.GetName() == "zai_quota_tool_calls_by_tool" {
			assert.Len(t, mf.GetMetric(), 0)
		}
	}
}

func TestMetrics_Update_InfoMetric(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{}
	metrics.Update(limits, "pro", time.Now(), time.Millisecond, true)

	expected := `# HELP zai_quota_info Information about the Z.ai subscription
# TYPE zai_quota_info gauge
zai_quota_info{level="pro"} 1
`
	err := testutil.CollectAndCompare(metrics.Info, strings.NewReader(expected), "zai_quota_info")
	assert.NoError(t, err)
}

func TestMetrics_Update_MultipleLimits(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits := []models.Limit{
		{
			Type:          "TOKENS_LIMIT",
			Percentage:    75,
			NextResetTime: 1710123400000,
		},
		{
			Type:          "TIME_LIMIT",
			CurrentValue:  50,
			Usage:         500,
			Remaining:     450,
			NextResetTime: 1710209800000,
		},
	}

	metrics.Update(limits, "enterprise", time.Now(), time.Millisecond, true)

	assert.Equal(t, 0.75, testutil.ToFloat64(metrics.PromptUsageRatio))
	assert.Equal(t, float64(50), testutil.ToFloat64(metrics.ToolCallsUsed))
	assert.Equal(t, float64(500), testutil.ToFloat64(metrics.ToolCallsLimit))
	assert.Equal(t, float64(450), testutil.ToFloat64(metrics.ToolCallsRemaining))
}

func TestMetrics_Update_ResetsToolCallsByTool(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := NewMetrics(reg)

	limits1 := []models.Limit{
		{
			Type: "TIME_LIMIT",
			UsageDetails: []models.UsageDetail{
				{ModelCode: "tool-a", Usage: 5},
				{ModelCode: "tool-b", Usage: 3},
			},
		},
	}

	metrics.Update(limits1, "pro", time.Now(), time.Millisecond, true)

	limits2 := []models.Limit{
		{
			Type: "TIME_LIMIT",
			UsageDetails: []models.UsageDetail{
				{ModelCode: "tool-c", Usage: 2},
			},
		},
	}

	metrics.Update(limits2, "pro", time.Now(), time.Millisecond, true)

	metricFamilies, err := reg.Gather()
	assert.NoError(t, err)

	for _, mf := range metricFamilies {
		if mf.GetName() == "zai_quota_tool_calls_by_tool" {
			assert.Len(t, mf.GetMetric(), 1, "Old tool metrics should be reset")
		}
	}
}
