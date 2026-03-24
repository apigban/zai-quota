package formatter

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"zai-quota/internal/processor"
)

func TestFormatJSON(t *testing.T) {
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

	result, err := FormatJSON(limits, "pro")
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	// Test 1: Output is valid JSON (parseable)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Test 2: JSON contains "limits" array
	limitsArray, ok := parsed["limits"]
	if !ok {
		t.Fatal("JSON does not contain 'limits' array")
	}

	limitsSlice, ok := limitsArray.([]interface{})
	if !ok {
		t.Fatal("'limits' is not an array")
	}

	if len(limitsSlice) != 2 {
		t.Fatalf("Expected 2 limits, got %d", len(limitsSlice))
	}

	// Test 3: Both limit types present
	limit0 := limitsSlice[0].(map[string]interface{})
	if limit0["type"] != "TOKENS_LIMIT" {
		t.Fatalf("First limit type should be TOKENS_LIMIT, got %v", limit0["type"])
	}

	limit1 := limitsSlice[1].(map[string]interface{})
	if limit1["type"] != "TIME_LIMIT" {
		t.Fatalf("Second limit type should be TIME_LIMIT, got %v", limit1["type"])
	}

	// Test 4: Check required fields in TOKENS_LIMIT (no used/total/remaining)
	requiredFieldsTokens := []string{"type", "label", "percentage", "next_reset", "next_reset_local"}
	for _, field := range requiredFieldsTokens {
		if _, ok := limit0[field]; !ok {
			t.Fatalf("Missing required field '%s' in TOKENS_LIMIT", field)
		}
	}

	// Test 5: Check required fields in TIME_LIMIT (includes used/total/remaining)
	requiredFieldsTime := []string{"type", "label", "percentage", "used", "total", "remaining", "next_reset", "next_reset_local"}
	for _, field := range requiredFieldsTime {
		if _, ok := limit1[field]; !ok {
			t.Fatalf("Missing required field '%s' in TIME_LIMIT", field)
		}
	}

	// Test 6: Verify field values for TIME_LIMIT
	if limit1["percentage"] != float64(25) {
		t.Fatalf("Expected percentage 25, got %v", limit1["percentage"])
	}
	if limit1["used"] != float64(1250) {
		t.Fatalf("Expected used 1250, got %v", limit1["used"])
	}
	if limit1["remaining"] != float64(3750) {
		t.Fatalf("Expected remaining 3750, got %v", limit1["remaining"])
	}

	// Test 7: Verify time formats
	nextReset, ok := limit0["next_reset"].(string)
	if !ok {
		t.Fatal("next_reset is not a string")
	}
	if !strings.Contains(nextReset, "2024-01-15T14:30:00") {
		t.Fatalf("Expected next_reset to contain '2024-01-15T14:30:00', got %s", nextReset)
	}

	nextResetLocal, ok := limit0["next_reset_local"].(string)
	if !ok {
		t.Fatal("next_reset_local is not a string")
	}
	if nextResetLocal != "2024-01-15 14:30" {
		t.Fatalf("Expected next_reset_local '2024-01-15 14:30', got %s", nextResetLocal)
	}

	// Test 8: Check indentation (pretty-printed with 2 spaces)
	if !strings.Contains(result, "  ") {
		t.Fatal("JSON should be pretty-printed with 2-space indentation")
	}
}

func TestFormatJSONEmptyArray(t *testing.T) {
	limits := []processor.ProcessedLimit{}

	_, err := FormatJSON(limits, "pro")
	if err == nil {
		t.Fatal("Expected error for empty limits array, got nil")
	}

	if !strings.Contains(err.Error(), "empty limits array") {
		t.Fatalf("Expected error message to contain 'empty limits array', got %v", err)
	}
}

func TestFormatJSONNilArray(t *testing.T) {
	_, err := FormatJSON(nil, "pro")
	if err == nil {
		t.Fatal("Expected error for nil limits array, got nil")
	}
}
