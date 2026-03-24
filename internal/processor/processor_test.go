package processor

import (
	"testing"
	"time"

	"zai-quota/internal/models"
)

func TestProcessLimits_BothLimitTypes(t *testing.T) {
	limits := []models.Limit{
		{
			Type:          "TOKENS_LIMIT",
			Unit:          3,
			Number:        5,
			Usage:         1000,
			CurrentValue:  250,
			Remaining:     750,
			Percentage:    25,
			NextResetTime: 1741401600000,
			UsageDetails: []models.UsageDetail{
				{ModelCode: "glm-4.7", Usage: 200},
			},
		},
		{
			Type:          "TIME_LIMIT",
			Usage:         500,
			CurrentValue:  100,
			Remaining:     400,
			Percentage:    20,
			NextResetTime: 1743993600000,
			UsageDetails: []models.UsageDetail{
				{ModelCode: "search-prime", Usage: 80},
			},
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 processed limits, got %d", len(result))
	}

	// Check TOKENS_LIMIT
	if result[0].Type != "TOKENS_LIMIT" {
		t.Errorf("Expected type TOKENS_LIMIT, got %s", result[0].Type)
	}
	if result[0].Label != "[5-Hour Prompt Limit]" {
		t.Errorf("Expected label '[5-Hour Prompt Limit]', got %s", result[0].Label)
	}
	if result[0].Percentage != 25 {
		t.Errorf("Expected percentage 25, got %d", result[0].Percentage)
	}
	if result[0].Used != 250 {
		t.Errorf("Expected used 250, got %d", result[0].Used)
	}
	if result[0].Total != 1000 {
		t.Errorf("Expected total 1000, got %d", result[0].Total)
	}
	if result[0].Remaining != 750 {
		t.Errorf("Expected remaining 750, got %d", result[0].Remaining)
	}
	if len(result[0].UsageDetails) != 1 {
		t.Errorf("Expected 1 usage detail, got %d", len(result[0].UsageDetails))
	}

	// Check TIME_LIMIT
	if result[1].Type != "TIME_LIMIT" {
		t.Errorf("Expected type TIME_LIMIT, got %s", result[1].Type)
	}
	if result[1].Label != "[Tool Quota]" {
		t.Errorf("Expected label '[Tool Quota]', got %s", result[1].Label)
	}
	if result[1].Percentage != 20 {
		t.Errorf("Expected percentage 20, got %d", result[1].Percentage)
	}
	if result[1].Used != 100 {
		t.Errorf("Expected used 100, got %d", result[1].Used)
	}
	if result[1].Total != 500 {
		t.Errorf("Expected total 500, got %d", result[1].Total)
	}
}

func TestProcessLimits_PercentageCapping(t *testing.T) {
	limits := []models.Limit{
		{
			Type:          "TIME_LIMIT",
			Percentage:    150, // Should be capped at 100
			Usage:         1000,
			CurrentValue:  800,
			Remaining:     200,
			NextResetTime: 1741401600000,
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 processed limit, got %d", len(result))
	}

	if result[0].Percentage != 100 {
		t.Errorf("Expected percentage to be capped at 100, got %d", result[0].Percentage)
	}
}

func TestProcessLimits_NegativeValueHandling(t *testing.T) {
	limits := []models.Limit{
		{
			Type:          "TIME_LIMIT",
			Percentage:    75,
			Usage:         -100, // Should be set to 0
			CurrentValue:  -50,  // Should be set to 0
			Remaining:     -25,  // Should be set to 0
			NextResetTime: 1741401600000,
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 processed limit, got %d", len(result))
	}

	if result[0].Used != 0 {
		t.Errorf("Expected used to be set to 0, got %d", result[0].Used)
	}

	if result[0].Total != 0 {
		t.Errorf("Expected total to be set to 0, got %d", result[0].Total)
	}

	if result[0].Remaining != 0 {
		t.Errorf("Expected remaining to be set to 0, got %d", result[0].Remaining)
	}
}

func TestProcessLimits_UnknownLimitType(t *testing.T) {
	limits := []models.Limit{
		{
			Type:          "TIME_LIMIT",
			Percentage:    75,
			Usage:         1000,
			CurrentValue:  250,
			Remaining:     750,
			NextResetTime: 1741401600000,
		},
		{
			Type:          "UNKNOWN_LIMIT", // Should be skipped with warning
			Percentage:    50,
			NextResetTime: 1743993600000,
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	// Unknown limit type should be skipped, so we expect only 1 result
	if len(result) != 1 {
		t.Fatalf("Expected 1 processed limit (unknown type skipped), got %d", len(result))
	}

	if result[0].Type != "TIME_LIMIT" {
		t.Errorf("Expected type TIME_LIMIT, got %s", result[0].Type)
	}
}

func TestProcessLimits_EmptyLimits(t *testing.T) {
	limits := []models.Limit{}

	result, err := ProcessLimits(limits)

	if err == nil {
		t.Fatalf("ProcessLimits should return error for empty limits array, got nil")
	}

	if result != nil {
		t.Fatalf("ProcessLimits should return nil result for empty limits, got %v", result)
	}

	expectedError := "empty limits array"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestProcessLimits_TimestampConversionAndFormatting(t *testing.T) {
	// Using a known timestamp: 1741401600000 = 2025-01-07 00:00:00 UTC
	limits := []models.Limit{
		{
			Type:          "TIME_LIMIT",
			Percentage:    75,
			Usage:         1000,
			CurrentValue:  750,
			Remaining:     250,
			NextResetTime: 1741401600000,
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 processed limit, got %d", len(result))
	}

	// Verify timestamp conversion
	expectedTime := time.Unix(1741401600000/1000, 0).Local()
	if !result[0].NextResetTime.Equal(expectedTime) {
		t.Errorf("Expected NextResetTime %v, got %v", expectedTime, result[0].NextResetTime)
	}

	// Verify formatting (format should be "2006-01-02 15:04")
	expectedFormat := expectedTime.Format("2006-01-02 15:04")
	if result[0].ResetFormatted != expectedFormat {
		t.Errorf("Expected ResetFormatted '%s', got '%s'", expectedFormat, result[0].ResetFormatted)
	}
}

func TestProcessLimits_FieldSemantics(t *testing.T) {
	// Test that both limits correctly map Usage -> Total, CurrentValue -> Used
	limits := []models.Limit{
		{
			Type:          "TOKENS_LIMIT",
			Usage:         1000,
			CurrentValue:  250,
			Remaining:     750,
			Percentage:    25,
			NextResetTime: 1741401600000,
		},
		{
			Type:          "TIME_LIMIT",
			Usage:         500,
			CurrentValue:  100,
			Remaining:     400,
			Percentage:    20,
			NextResetTime: 1743993600000,
		},
	}

	result, err := ProcessLimits(limits)

	if err != nil {
		t.Fatalf("ProcessLimits returned error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 processed limits, got %d", len(result))
	}

	// TOKENS_LIMIT
	if result[0].Total != 1000 {
		t.Errorf("Expected Total 1000 for TOKENS_LIMIT, got %d", result[0].Total)
	}
	if result[0].Used != 250 {
		t.Errorf("Expected Used 250 for TOKENS_LIMIT, got %d", result[0].Used)
	}

	// TIME_LIMIT
	if result[1].Total != 500 {
		t.Errorf("Expected Total 500 for TIME_LIMIT, got %d", result[1].Total)
	}
	if result[1].Used != 100 {
		t.Errorf("Expected Used 100 for TIME_LIMIT, got %d", result[1].Used)
	}
}
