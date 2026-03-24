package processor

import (
	"fmt"
	"time"
)

// ConvertTimestamp converts a millisecond timestamp to local time
func ConvertTimestamp(ms int64) time.Time {
	return time.Unix(ms/1000, 0).Local()
}

// FormatTimeUntil formats the duration until a future time
func FormatTimeUntil(t time.Time) string {
	d := time.Until(t)
	if d < 0 {
		return "Past"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		return fmt.Sprintf("%d days", days)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	return fmt.Sprintf("%dm", minutes)
}

func FormatResetDateTime(t time.Time) string {
	iana, found := GetTimezone()
	name, offsetSeconds := t.Zone()

	dateStr := t.Format("Jan 2, 15:04")

	if !found {
		if name == "" {
			name = "UTC"
		}
		return fmt.Sprintf("%s %s %s", dateStr, FormatTimezoneOffset(offsetSeconds), name)
	}

	offset := FormatTimezoneOffset(offsetSeconds)
	city := ExtractCity(iana)

	return fmt.Sprintf("%s %s %s", dateStr, offset, city)
}
