package formatter

import (
	"strings"
	"testing"
	"time"

	"zai-quota/internal/processor"
)

func TestFormatYAML(t *testing.T) {
	// Create sample limits
	tokensTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	timeTime := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	limits := []processor.ProcessedLimit{
		{
			Type:           "TOKENS_LIMIT",
			Label:          "[5-Hour Prompt Limit]",
			Percentage:     75,
			Used:           0,
			Total:          0,
			Remaining:      0,
			NextResetTime:  tokensTime,
			ResetFormatted: "2024-01-15 14:30",
		},
		{
			Type:           "TIME_LIMIT",
			Label:          "[Tool Quota]",
			Percentage:     25,
			Used:           1250,
			Total:          5000,
			Remaining:      3750,
			NextResetTime:  timeTime,
			ResetFormatted: "2024-02-01 00:00",
		},
	}

	result, err := FormatYAML(limits, "pro")
	if err != nil {
		t.Fatalf("FormatYAML failed: %v", err)
	}

	// Test 1: Output contains "limits" key
	if !strings.Contains(result, "limits:") {
		t.Fatal("YAML does not contain 'limits' key")
	}

	// Test 2: Both limit types present
	if !strings.Contains(result, "type: TOKENS_LIMIT") {
		t.Fatal("YAML does not contain TOKENS_LIMIT type")
	}
	if !strings.Contains(result, "type: TIME_LIMIT") {
		t.Fatal("YAML does not contain TIME_LIMIT type")
	}

	// Test 3: Check required fields (common to both types)
	requiredFields := []string{"type:", "label:", "percentage:", "next_reset:", "next_reset_local:"}
	for _, field := range requiredFields {
		if !strings.Contains(result, field) {
			t.Fatalf("Missing required field '%s' in YAML output", field)
		}
	}

	// Test 4: Verify field values for TIME_LIMIT
	if !strings.Contains(result, "percentage: 25") {
		t.Fatal("Expected percentage 25 not found in YAML")
	}
	if !strings.Contains(result, "used: 1250") {
		t.Fatal("Expected used 1250 not found in YAML")
	}
	if !strings.Contains(result, "remaining: 3750") {
		t.Fatal("Expected remaining 3750 not found in YAML")
	}

	// Test 5: Verify time formats
	if !strings.Contains(result, "2024-01-15T14:30:00") {
		t.Fatalf("Expected next_reset to contain '2024-01-15T14:30:00', got: %s", result)
	}
	if !strings.Contains(result, "2024-01-15 14:30") {
		t.Fatalf("Expected next_reset_local '2024-01-15 14:30', got: %s", result)
	}

	// Test 6: Check indentation (YAML with 2 spaces)
	if !strings.Contains(result, "  ") {
		t.Fatal("YAML should use 2-space indentation")
	}
}

func TestFormatYAMLEmptyArray(t *testing.T) {
	limits := []processor.ProcessedLimit{}

	_, err := FormatYAML(limits, "pro")
	if err == nil {
		t.Fatal("Expected error for empty limits array, got nil")
	}

	if !strings.Contains(err.Error(), "empty limits array") {
		t.Fatalf("Expected error message to contain 'empty limits array', got %v", err)
	}
}

func TestFormatYAMLNilArray(t *testing.T) {
	_, err := FormatYAML(nil, "pro")
	if err == nil {
		t.Fatal("Expected error for nil limits array, got nil")
	}
}

func TestFormatYAMLSpecialChars(t *testing.T) {
	// Test with special characters in type name
	specialTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	limits := []processor.ProcessedLimit{
		{
			Type:           "LIMIT_WITH_SPECIAL_CHARS!@#$",
			Percentage:     50,
			Used:           5000,
			Remaining:      5000,
			NextResetTime:  specialTime,
			ResetFormatted: "2024-01-15 14:30",
		},
	}

	result, err := FormatYAML(limits, "pro")
	if err != nil {
		t.Fatalf("FormatYAML failed with special chars: %v", err)
	}

	if !strings.Contains(result, "LIMIT_WITH_SPECIAL_CHARS") {
		t.Fatal("YAML should contain special characters in type")
	}
}

func TestFormatYAMLUnicode(t *testing.T) {
	// Test with unicode characters
	unicodeTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	limits := []processor.ProcessedLimit{
		{
			Type:           "LIMIT_中文_测试",
			Percentage:     75,
			Used:           7500,
			Remaining:      2500,
			NextResetTime:  unicodeTime,
			ResetFormatted: "2024-01-15 14:30",
		},
	}

	result, err := FormatYAML(limits, "pro")
	if err != nil {
		t.Fatalf("FormatYAML failed with unicode: %v", err)
	}

	if !strings.Contains(result, "LIMIT_中文_测试") {
		t.Fatal("YAML should contain unicode characters in type")
	}
}

func TestFormatYAMLZeroValues(t *testing.T) {
	// Test with zero percentage and usage values
	zeroTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	limits := []processor.ProcessedLimit{
		{
			Type:           "ZERO_LIMIT",
			Percentage:     0,
			Used:           0,
			Remaining:      10000,
			NextResetTime:  zeroTime,
			ResetFormatted: "2024-01-15 14:30",
		},
	}

	result, err := FormatYAML(limits, "pro")
	if err != nil {
		t.Fatalf("FormatYAML failed with zero values: %v", err)
	}

	if !strings.Contains(result, "percentage: 0") {
		t.Fatal("YAML should contain percentage 0")
	}
}
