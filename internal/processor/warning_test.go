package processor

import (
	"testing"
)

func TestCalculateWarningLevel(t *testing.T) {
	tests := []struct {
		name       string
		percentage int
		expected   string
	}{
		// Safe level (< 80)
		{
			name:       "79 percent should be safe",
			percentage: 79,
			expected:   "safe",
		},
		{
			name:       "50 percent should be safe",
			percentage: 50,
			expected:   "safe",
		},
		{
			name:       "0 percent should be safe",
			percentage: 0,
			expected:   "safe",
		},

		// Warning level (80-89)
		{
			name:       "80 percent should be warning",
			percentage: 80,
			expected:   "warning",
		},
		{
			name:       "81 percent should be warning",
			percentage: 81,
			expected:   "warning",
		},
		{
			name:       "89 percent should be warning",
			percentage: 89,
			expected:   "warning",
		},

		// Critical level (90-94)
		{
			name:       "90 percent should be critical",
			percentage: 90,
			expected:   "critical",
		},
		{
			name:       "91 percent should be critical",
			percentage: 91,
			expected:   "critical",
		},
		{
			name:       "94 percent should be critical",
			percentage: 94,
			expected:   "critical",
		},

		// Emergency level (>= 95)
		{
			name:       "95 percent should be emergency",
			percentage: 95,
			expected:   "emergency",
		},
		{
			name:       "96 percent should be emergency",
			percentage: 96,
			expected:   "emergency",
		},
		{
			name:       "100 percent should be emergency",
			percentage: 100,
			expected:   "emergency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWarningLevel(tt.percentage)

			if result != tt.expected {
				t.Errorf("CalculateWarningLevel(%d) = %s, want %s", tt.percentage, result, tt.expected)
			}
		})
	}
}

func TestCalculateWarningLevel_BoundaryConditions(t *testing.T) {
	// Explicit boundary testing to ensure correct transitions
	boundaryTests := []struct {
		name       string
		percentage int
		expected   string
	}{
		// Just below warning threshold
		{"79 (below warning threshold)", 79, "safe"},

		// At warning threshold
		{"80 (at warning threshold)", 80, "warning"},

		// In warning range
		{"81 (in warning range)", 81, "warning"},

		// Just below critical threshold
		{"89 (below critical threshold)", 89, "warning"},

		// At critical threshold
		{"90 (at critical threshold)", 90, "critical"},

		// In critical range
		{"91 (in critical range)", 91, "critical"},

		// Just below emergency threshold
		{"94 (below emergency threshold)", 94, "critical"},

		// At emergency threshold
		{"95 (at emergency threshold)", 95, "emergency"},

		// In emergency range
		{"96 (in emergency range)", 96, "emergency"},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWarningLevel(tt.percentage)

			if result != tt.expected {
				t.Errorf("CalculateWarningLevel(%d) = %s, want %s", tt.percentage, result, tt.expected)
			}
		})
	}
}
