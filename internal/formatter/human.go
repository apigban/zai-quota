package formatter

import (
	"fmt"
	"strings"

	"zai-quota/internal/processor"
)

func FormatHuman(limits []processor.ProcessedLimit) string {
	var output strings.Builder

	for _, limit := range limits {
		if limit.Type != "TOKENS_LIMIT" && limit.Type != "TIME_LIMIT" {
			continue
		}
		fmt.Fprintf(&output, "%s\n", limit.Label)
		fmt.Fprintf(&output, "Type: %s\n", limit.Type)

		if limit.Type == "TIME_LIMIT" {
			fmt.Fprintf(&output, "Usage: %d / %d\n", limit.Used, limit.Total)
			fmt.Fprintf(&output, "Remaining: %d\n", limit.Remaining)
			if len(limit.UsageDetails) > 0 {
				fmt.Fprintf(&output, "Usage breakdown:\n")
				for _, detail := range limit.UsageDetails {
					fmt.Fprintf(&output, "  └─ %s: %d\n", detail.ModelCode, detail.Usage)
				}
			}
		} else if limit.Type == "TOKENS_LIMIT" {
			fmt.Fprintf(&output, "Usage: %d%%\n", limit.Percentage)
		}

		fmt.Fprintf(&output, "Next Reset: %s\n", limit.ResetFormatted)
		fmt.Fprintln(&output)
	}

	return output.String()
}
