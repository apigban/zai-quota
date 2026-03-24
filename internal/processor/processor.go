package processor

import (
	"fmt"
	"log"
	"time"

	"zai-quota/internal/models"
)

type ProcessedLimit struct {
	Type           string
	Label          string
	Percentage     int
	Total          int
	Used           int
	Remaining      int
	NextResetTime  time.Time
	ResetFormatted string
	UsageDetails   []models.UsageDetail
}

func ProcessLimits(limits []models.Limit) ([]ProcessedLimit, error) {
	if len(limits) == 0 {
		return nil, fmt.Errorf("empty limits array")
	}

	processed := make([]ProcessedLimit, 0, len(limits))

	for _, limit := range limits {
		if limit.Type != "TOKENS_LIMIT" && limit.Type != "TIME_LIMIT" {
			log.Printf("WARNING: unknown limit type '%s', skipping\n", limit.Type)
			continue
		}

		percentage := limit.Percentage
		if percentage > 100 {
			percentage = 100
		}
		if percentage < 0 {
			percentage = 0
		}

		var label string
		var used, total, remaining int
		var usageDetails []models.UsageDetail

		if limit.Type == "TOKENS_LIMIT" {
			label = "[5-Hour Prompt Limit]"
			used = limit.CurrentValue
			if used < 0 {
				used = 0
			}
			total = limit.Usage
			if total < 0 {
				total = 0
			}
			remaining = limit.Remaining
			if remaining < 0 {
				remaining = 0
			}
			usageDetails = limit.UsageDetails
		} else if limit.Type == "TIME_LIMIT" {
			label = "[Tool Quota]"
			used = limit.CurrentValue
			if used < 0 {
				used = 0
			}
			total = limit.Usage
			if total < 0 {
				total = 0
			}
			remaining = limit.Remaining
			if remaining < 0 {
				remaining = 0
			}
			usageDetails = limit.UsageDetails
		}

		resetTime := ConvertTimestamp(limit.NextResetTime)
		resetFormatted := resetTime.Format("2006-01-02 15:04")

		processed = append(processed, ProcessedLimit{
			Type:           limit.Type,
			Label:          label,
			Percentage:     percentage,
			Total:          total,
			Used:           used,
			Remaining:      remaining,
			NextResetTime:  resetTime,
			ResetFormatted: resetFormatted,
			UsageDetails:   usageDetails,
		})
	}

	return processed, nil
}
