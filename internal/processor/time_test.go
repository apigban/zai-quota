package processor

import (
	"testing"
	"time"
)

func TestConvertTimestamp(t *testing.T) {
	// Test standard millisecond conversion
	t.Run("millisecond conversion", func(t *testing.T) {
		// 2024-01-01 12:00:00 UTC in milliseconds
		ms := int64(1704110400000)
		result := ConvertTimestamp(ms)

		// Verify the time is approximately correct (within 1 second)
		expected := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		diff := result.Sub(expected)
		if diff < -time.Second || diff > time.Second {
			t.Errorf("ConvertTimestamp(%d) = %v, expected approximately %v", ms, result, expected)
		}
	})

	// Test local timezone conversion
	t.Run("local timezone conversion", func(t *testing.T) {
		ms := int64(1704110400000)
		result := ConvertTimestamp(ms)

		// Verify result is in local timezone
		if result.Location().String() == "UTC" {
			t.Errorf("ConvertTimestamp should return local time, got UTC")
		}
	})

	// Test zero timestamp handling
	t.Run("zero timestamp", func(t *testing.T) {
		ms := int64(0)
		result := ConvertTimestamp(ms)

		// Zero timestamp should not panic and should return Unix epoch
		expected := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		diff := result.Sub(expected)
		if diff < -time.Second || diff > time.Second {
			t.Errorf("ConvertTimestamp(%d) = %v, expected approximately %v", ms, result, expected)
		}
	})

	// Test negative timestamp handling
	t.Run("negative timestamp", func(t *testing.T) {
		// 1969-12-31 23:59:59 UTC in milliseconds
		ms := int64(-1000)
		result := ConvertTimestamp(ms)

		// Should handle negative timestamps gracefully
		expected := time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC)
		diff := result.Sub(expected)
		if diff < -time.Second || diff > time.Second {
			t.Errorf("ConvertTimestamp(%d) = %v, expected approximately %v", ms, result, expected)
		}
	})

	// Test current time conversion
	t.Run("current time", func(t *testing.T) {
		now := time.Now()
		ms := now.UnixMilli()
		result := ConvertTimestamp(ms)

		// Result should be within 1 second of original time
		diff := result.Sub(now)
		if diff < -time.Second || diff > time.Second {
			t.Errorf("ConvertTimestamp(%d) = %v, expected approximately %v", ms, result, now)
		}
	})
}

func TestFormatResetDateTime(t *testing.T) {
	t.Run("returns formatted string with expected components", func(t *testing.T) {
		testTime := time.Date(2024, 3, 15, 14, 30, 0, 0, time.Local)
		result := FormatResetDateTime(testTime)

		if result == "" {
			t.Error("FormatResetDateTime returned empty string")
		}

		if len(result) < 10 {
			t.Errorf("FormatResetDateTime result too short: %q", result)
		}
	})

	t.Run("contains date and time", func(t *testing.T) {
		testTime := time.Date(2024, 3, 15, 14, 30, 0, 0, time.Local)
		result := FormatResetDateTime(testTime)

		if result == "UTC" || result == "" {
			t.Skip("Timezone detection not available on this system")
		}

		if result[0] < 'A' || result[0] > 'Z' {
			t.Errorf("Expected date to start with month abbreviation, got: %q", result)
		}
	})

	t.Run("handles UTC fallback gracefully", func(t *testing.T) {
		testTime := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
		result := FormatResetDateTime(testTime)

		if result == "" {
			t.Error("FormatResetDateTime should not return empty string")
		}
	})
}
