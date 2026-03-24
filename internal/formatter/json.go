package formatter

import (
	"encoding/json"
	"fmt"
	"time"

	"zai-quota/internal/processor"
)

type JSONLimit struct {
	Type           string        `json:"type"`
	Label          string        `json:"label"`
	Percentage     int           `json:"percentage"`
	Used           int           `json:"used,omitempty"`
	Total          int           `json:"total,omitempty"`
	Remaining      int           `json:"remaining,omitempty"`
	NextReset      string        `json:"next_reset"`
	NextResetLocal string        `json:"next_reset_local"`
	UsageDetails   []UsageDetail `json:"usage_details,omitempty"`
}

type UsageDetail struct {
	ModelCode string `json:"model_code"`
	Usage     int    `json:"usage"`
}

type JSONOutput struct {
	Limits []JSONLimit `json:"limits"`
	Level  string      `json:"level"`
}

func FormatJSON(limits []processor.ProcessedLimit, level string) (string, error) {
	if len(limits) == 0 {
		return "", fmt.Errorf("empty limits array")
	}

	jsonLimits := make([]JSONLimit, 0, len(limits))

	for _, limit := range limits {
		jsonLimit := JSONLimit{
			Type:           limit.Type,
			Label:          limit.Label,
			Percentage:     limit.Percentage,
			NextReset:      limit.NextResetTime.Format(time.RFC3339),
			NextResetLocal: limit.ResetFormatted,
		}

		if limit.Type == "TIME_LIMIT" {
			jsonLimit.Used = limit.Used
			jsonLimit.Total = limit.Total
			jsonLimit.Remaining = limit.Remaining
			if len(limit.UsageDetails) > 0 {
				jsonLimit.UsageDetails = make([]UsageDetail, 0, len(limit.UsageDetails))
				for _, detail := range limit.UsageDetails {
					jsonLimit.UsageDetails = append(jsonLimit.UsageDetails, UsageDetail{
						ModelCode: detail.ModelCode,
						Usage:     detail.Usage,
					})
				}
			}
		}

		jsonLimits = append(jsonLimits, jsonLimit)
	}

	output := JSONOutput{
		Limits: jsonLimits,
		Level:  level,
	}

	result, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(result), nil
}
