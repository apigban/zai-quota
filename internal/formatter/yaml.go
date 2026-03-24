package formatter

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
	"zai-quota/internal/processor"
)

type YAMLLimit struct {
	Type           string        `yaml:"type"`
	Label          string        `yaml:"label"`
	Percentage     int           `yaml:"percentage"`
	Used           int           `yaml:"used,omitempty"`
	Total          int           `yaml:"total,omitempty"`
	Remaining      int           `yaml:"remaining,omitempty"`
	NextReset      string        `yaml:"next_reset"`
	NextResetLocal string        `yaml:"next_reset_local"`
	UsageDetails   []UsageDetail `yaml:"usage_details,omitempty"`
}

type YAMLOutput struct {
	Limits []YAMLLimit `yaml:"limits"`
	Level  string      `yaml:"level"`
}

func FormatYAML(limits []processor.ProcessedLimit, level string) (string, error) {
	if len(limits) == 0 {
		return "", fmt.Errorf("empty limits array")
	}

	yamlLimits := make([]YAMLLimit, 0, len(limits))

	for _, limit := range limits {
		yamlLimit := YAMLLimit{
			Type:           limit.Type,
			Label:          limit.Label,
			Percentage:     limit.Percentage,
			NextReset:      limit.NextResetTime.Format(time.RFC3339),
			NextResetLocal: limit.ResetFormatted,
		}

		if limit.Type == "TIME_LIMIT" {
			yamlLimit.Used = limit.Used
			yamlLimit.Total = limit.Total
			yamlLimit.Remaining = limit.Remaining
			if len(limit.UsageDetails) > 0 {
				yamlLimit.UsageDetails = make([]UsageDetail, 0, len(limit.UsageDetails))
				for _, detail := range limit.UsageDetails {
					yamlLimit.UsageDetails = append(yamlLimit.UsageDetails, UsageDetail{
						ModelCode: detail.ModelCode,
						Usage:     detail.Usage,
					})
				}
			}
		}

		yamlLimits = append(yamlLimits, yamlLimit)
	}

	output := YAMLOutput{
		Limits: yamlLimits,
		Level:  level,
	}

	result, err := yaml.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(result), nil
}
