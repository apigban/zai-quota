package models

import (
	"encoding/json"
	"testing"
)

func TestQuotaResponse_UnmarshalJSON_Valid(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     QuotaResponse
	}{
		{
			name: "Full API response structure",
			jsonData: `{
				"limits": [
					{
						"type": "TIME_LIMIT",
						"unit": 5,
						"number": 1,
						"usage": 1000,
						"currentValue": 8,
						"remaining": 992,
						"percentage": 1,
						"nextResetTime": 1775186469998,
						"usageDetails": [
							{"modelCode": "search-prime", "usage": 5},
							{"modelCode": "model-alpha", "usage": 3}
						]
					},
					{
						"type": "TOKENS_LIMIT",
						"unit": 3,
						"number": 5,
						"percentage": 33,
						"nextResetTime": 1773052446291
					}
				],
				"level": "pro"
			}`,
			want: QuotaResponse{
				Limits: []Limit{
					{
						Type:          "TIME_LIMIT",
						Unit:          5,
						Number:        1,
						Usage:         1000,
						CurrentValue:  8,
						Remaining:     992,
						Percentage:    1,
						NextResetTime: 1775186469998,
						UsageDetails: []UsageDetail{
							{ModelCode: "search-prime", Usage: 5},
							{ModelCode: "model-alpha", Usage: 3},
						},
					},
					{
						Type:          "TOKENS_LIMIT",
						Unit:          3,
						Number:        5,
						Percentage:    33,
						NextResetTime: 1773052446291,
					},
				},
				Level: "pro",
			},
		},
		{
			name: "Empty limits",
			jsonData: `{
				"limits": [],
				"level": "lite"
			}`,
			want: QuotaResponse{
				Limits: []Limit{},
				Level:  "lite",
			},
		},
		{
			name: "TIME_LIMIT only",
			jsonData: `{
				"limits": [
					{
						"type": "TIME_LIMIT",
						"unit": 5,
						"number": 1,
						"usage": 1000,
						"currentValue": 8,
						"remaining": 992,
						"percentage": 1,
						"nextResetTime": 1775186469998
					}
				],
				"level": "max"
			}`,
			want: QuotaResponse{
				Limits: []Limit{
					{
						Type:          "TIME_LIMIT",
						Unit:          5,
						Number:        1,
						Usage:         1000,
						CurrentValue:  8,
						Remaining:     992,
						Percentage:    1,
						NextResetTime: 1775186469998,
					},
				},
				Level: "max",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got QuotaResponse
			err := json.Unmarshal([]byte(tt.jsonData), &got)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}
			if got.Level != tt.want.Level {
				t.Errorf("Level mismatch: got %q, want %q", got.Level, tt.want.Level)
			}
			if len(got.Limits) != len(tt.want.Limits) {
				t.Fatalf("Limits count mismatch: got %d, want %d", len(got.Limits), len(tt.want.Limits))
			}
			for i := range got.Limits {
				if got.Limits[i].Type != tt.want.Limits[i].Type {
					t.Errorf("Limit[%d].Type mismatch: got %q, want %q", i, got.Limits[i].Type, tt.want.Limits[i].Type)
				}
			}
		})
	}
}

func TestLimit_UnmarshalJSON_TIME_LIMIT(t *testing.T) {
	jsonData := `{
		"type": "TIME_LIMIT",
		"unit": 5,
		"number": 1,
		"usage": 1000,
		"currentValue": 8,
		"remaining": 992,
		"percentage": 1,
		"nextResetTime": 1775186469998,
		"usageDetails": [
			{"modelCode": "search-prime", "usage": 5},
			{"modelCode": "model-alpha", "usage": 3}
		]
	}`

	var got Limit
	err := json.Unmarshal([]byte(jsonData), &got)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if got.Type != "TIME_LIMIT" {
		t.Errorf("Type mismatch: got %q, want %q", got.Type, "TIME_LIMIT")
	}
	if got.Unit != 5 {
		t.Errorf("Unit mismatch: got %d, want %d", got.Unit, 5)
	}
	if got.Number != 1 {
		t.Errorf("Number mismatch: got %d, want %d", got.Number, 1)
	}
	if got.Usage != 1000 {
		t.Errorf("Usage mismatch: got %d, want %d", got.Usage, 1000)
	}
	if got.CurrentValue != 8 {
		t.Errorf("CurrentValue mismatch: got %d, want %d", got.CurrentValue, 8)
	}
	if got.Remaining != 992 {
		t.Errorf("Remaining mismatch: got %d, want %d", got.Remaining, 992)
	}
	if got.Percentage != 1 {
		t.Errorf("Percentage mismatch: got %d, want %d", got.Percentage, 1)
	}
	if got.NextResetTime != 1775186469998 {
		t.Errorf("NextResetTime mismatch: got %d, want %d", got.NextResetTime, 1775186469998)
	}
	if len(got.UsageDetails) != 2 {
		t.Fatalf("UsageDetails count mismatch: got %d, want %d", len(got.UsageDetails), 2)
	}
	if got.UsageDetails[0].ModelCode != "search-prime" {
		t.Errorf("UsageDetails[0].ModelCode mismatch: got %q, want %q", got.UsageDetails[0].ModelCode, "search-prime")
	}
	if got.UsageDetails[0].Usage != 5 {
		t.Errorf("UsageDetails[0].Usage mismatch: got %d, want %d", got.UsageDetails[0].Usage, 5)
	}
}

func TestLimit_UnmarshalJSON_TOKENS_LIMIT(t *testing.T) {
	jsonData := `{
		"type": "TOKENS_LIMIT",
		"unit": 3,
		"number": 5,
		"percentage": 33,
		"nextResetTime": 1773052446291
	}`

	var got Limit
	err := json.Unmarshal([]byte(jsonData), &got)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if got.Type != "TOKENS_LIMIT" {
		t.Errorf("Type mismatch: got %q, want %q", got.Type, "TOKENS_LIMIT")
	}
	if got.Unit != 3 {
		t.Errorf("Unit mismatch: got %d, want %d", got.Unit, 3)
	}
	if got.Number != 5 {
		t.Errorf("Number mismatch: got %d, want %d", got.Number, 5)
	}
	if got.Percentage != 33 {
		t.Errorf("Percentage mismatch: got %d, want %d", got.Percentage, 33)
	}
	if got.NextResetTime != 1773052446291 {
		t.Errorf("NextResetTime mismatch: got %d, want %d", got.NextResetTime, 1773052446291)
	}
	// TOKENS_LIMIT should not have these fields
	if got.Usage != 0 {
		t.Errorf("Usage should be 0 for TOKENS_LIMIT, got %d", got.Usage)
	}
	if got.CurrentValue != 0 {
		t.Errorf("CurrentValue should be 0 for TOKENS_LIMIT, got %d", got.CurrentValue)
	}
	if got.Remaining != 0 {
		t.Errorf("Remaining should be 0 for TOKENS_LIMIT, got %d", got.Remaining)
	}
}

func TestLimit_UnmarshalJSON_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name:     "Invalid JSON syntax",
			jsonData: `{"type": "TIME_LIMIT", "percentage": 1`,
			wantErr:  true,
		},
		{
			name:     "Wrong type for type field",
			jsonData: `{"type": 123, "percentage": 1}`,
			wantErr:  true,
		},
		{
			name:     "Wrong type for percentage",
			jsonData: `{"type": "TIME_LIMIT", "percentage": "not-a-number"}`,
			wantErr:  true,
		},
		{
			name:     "Wrong type for nextResetTime",
			jsonData: `{"type": "TIME_LIMIT", "nextResetTime": "not-a-number"}`,
			wantErr:  true,
		},
		{
			name:     "Extra fields ignored",
			jsonData: `{"type": "TIME_LIMIT", "percentage": 1, "extra": "field"}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Limit
			err := json.Unmarshal([]byte(tt.jsonData), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsageDetail_UnmarshalJSON(t *testing.T) {
	jsonData := `{"modelCode": "search-prime", "usage": 5}`

	var got UsageDetail
	err := json.Unmarshal([]byte(jsonData), &got)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if got.ModelCode != "search-prime" {
		t.Errorf("ModelCode mismatch: got %q, want %q", got.ModelCode, "search-prime")
	}
	if got.Usage != 5 {
		t.Errorf("Usage mismatch: got %d, want %d", got.Usage, 5)
	}
}

func TestQuotaResponse_MarshalJSON(t *testing.T) {
	response := QuotaResponse{
		Limits: []Limit{
			{
				Type:          "TIME_LIMIT",
				Unit:          5,
				Number:        1,
				Usage:         1000,
				CurrentValue:  8,
				Remaining:     992,
				Percentage:    1,
				NextResetTime: 1775186469998,
				UsageDetails: []UsageDetail{
					{ModelCode: "search-prime", Usage: 5},
				},
			},
		},
		Level: "pro",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify it contains expected fields
	jsonStr := string(data)
	if !contains(jsonStr, `"type":"TIME_LIMIT"`) {
		t.Error("MarshalJSON should contain type field")
	}
	if !contains(jsonStr, `"level":"pro"`) {
		t.Error("MarshalJSON should contain level field")
	}
	if !contains(jsonStr, `"currentValue":8`) {
		t.Error("MarshalJSON should contain currentValue field")
	}
}

func TestLimit_MarshalJSON_OmitEmpty(t *testing.T) {
	// TOKENS_LIMIT should omit usage, currentValue, remaining, usageDetails
	limit := Limit{
		Type:          "TOKENS_LIMIT",
		Unit:          3,
		Number:        5,
		Percentage:    33,
		NextResetTime: 1773052446291,
	}

	data, err := json.Marshal(limit)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonStr := string(data)
	if contains(jsonStr, `"usage"`) {
		t.Error("MarshalJSON should omit usage field when zero for TOKENS_LIMIT")
	}
	if contains(jsonStr, `"currentValue"`) {
		t.Error("MarshalJSON should omit currentValue field when zero for TOKENS_LIMIT")
	}
	if contains(jsonStr, `"remaining"`) {
		t.Error("MarshalJSON should omit remaining field when zero for TOKENS_LIMIT")
	}
	if contains(jsonStr, `"usageDetails"`) {
		t.Error("MarshalJSON should omit usageDetails field when nil")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
