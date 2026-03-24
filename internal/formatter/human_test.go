package formatter

import (
	"strings"
	"testing"
	"time"

	"zai-quota/internal/processor"
)

func TestFormatHuman(t *testing.T) {
	tests := []struct {
		name     string
		limits   []processor.ProcessedLimit
		contains []string
	}{
		{
			name: "both limits present",
			limits: []processor.ProcessedLimit{
				{
					Type:           "TOKENS_LIMIT",
					Label:          "[5-Hour Prompt Limit]",
					Percentage:     75,
					Used:           0,
					Total:          0,
					Remaining:      0,
					NextResetTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
					ResetFormatted: "2024-01-15 14:30",
				},
				{
					Type:           "TIME_LIMIT",
					Label:          "[Tool Quota]",
					Percentage:     25,
					Used:           1250,
					Total:          5000,
					Remaining:      3750,
					NextResetTime:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					ResetFormatted: "2024-02-01 00:00",
				},
			},
			contains: []string{
				"[5-Hour Prompt Limit]",
				"[Tool Quota]",
				"Type: TOKENS_LIMIT",
				"Type: TIME_LIMIT",
				"Usage: 75%",
				"Usage: 1250 / 5000",
				"Remaining: 3750",
				"Next Reset: 2024-01-15 14:30",
				"Next Reset: 2024-02-01 00:00",
			},
		},
		{
			name: "100% usage edge case",
			limits: []processor.ProcessedLimit{
				{
					Type:           "TOKENS_LIMIT",
					Label:          "[5-Hour Prompt Limit]",
					Percentage:     100,
					Used:           0,
					Total:          0,
					Remaining:      0,
					NextResetTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
					ResetFormatted: "2024-01-15 14:30",
				},
			},
			contains: []string{
				"[5-Hour Prompt Limit]",
				"Type: TOKENS_LIMIT",
				"Usage: 100%",
				"Next Reset: 2024-01-15 14:30",
			},
		},
		{
			name: "0 remaining edge case",
			limits: []processor.ProcessedLimit{
				{
					Type:           "TIME_LIMIT",
					Label:          "[Tool Quota]",
					Percentage:     100,
					Used:           5000,
					Total:          5000,
					Remaining:      0,
					NextResetTime:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					ResetFormatted: "2024-02-01 00:00",
				},
			},
			contains: []string{
				"[Tool Quota]",
				"Type: TIME_LIMIT",
				"Usage: 5000 / 5000",
				"Remaining: 0",
				"Next Reset: 2024-02-01 00:00",
			},
		},
		{
			name:   "empty limits",
			limits: []processor.ProcessedLimit{},
			contains: []string{
				"", // Empty output
			},
		},
		{
			name: "unknown limit types are skipped",
			limits: []processor.ProcessedLimit{
				{
					Type:           "UNKNOWN_LIMIT",
					Percentage:     50,
					Used:           1000,
					Remaining:      1000,
					NextResetTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
					ResetFormatted: "2024-01-15 14:30",
				},
			},
			contains: []string{
				"", // Empty output (skipped unknown type)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FormatHuman(tt.limits)

			// If expecting empty output
			if len(tt.contains) == 1 && tt.contains[0] == "" {
				if output != "" {
					t.Errorf("Expected empty output, got: %s", output)
				}
				return
			}

			// Check that all expected strings are in output
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Output missing expected string:\nExpected: %s\nGot:\n%s", expected, output)
				}
			}

			// Verify no ANSI color codes
			if strings.Contains(output, "\x1b[") || strings.Contains(output, "\033[") {
				t.Errorf("Output contains ANSI color codes:\n%s", output)
			}

			// Verify no tabs
			if strings.Contains(output, "\t") {
				t.Errorf("Output contains tab characters:\n%s", output)
			}
		})
	}
}

func TestFormatHumanTemplateMatch(t *testing.T) {
	limits := []processor.ProcessedLimit{
		{
			Type:           "TOKENS_LIMIT",
			Label:          "[5-Hour Prompt Limit]",
			Percentage:     75,
			Used:           75000,
			Total:          100000,
			Remaining:      25000,
			NextResetTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			ResetFormatted: "2024-01-15 14:30",
		},
		{
			Type:           "TIME_LIMIT",
			Label:          "[Tool Quota]",
			Percentage:     25,
			Used:           1250,
			Total:          5000,
			Remaining:      3750,
			NextResetTime:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			ResetFormatted: "2024-02-01 00:00",
		},
	}

	output := FormatHuman(limits)

	// Check exact template structure
	expectedLines := []string{
		"[5-Hour Prompt Limit]",
		"Type: TOKENS_LIMIT",
		"Usage: 75%",
		"Next Reset: 2024-01-15 14:30",
		"",
		"[Tool Quota]",
		"Type: TIME_LIMIT",
		"Usage: 1250 / 5000",
		"Remaining: 3750",
		"Next Reset: 2024-02-01 00:00",
		"",
	}

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != len(expectedLines) {
		t.Errorf("Line count mismatch: got %d, want %d\nOutput:\n%s", len(lines), len(expectedLines), output)
	}

	for i, expected := range expectedLines {
		if i < len(lines) && lines[i] != expected {
			t.Errorf("Line %d mismatch:\nExpected: %s\nGot: %s\nFull output:\n%s", i, expected, lines[i], output)
		}
	}
}
