package processor

import (
	"testing"
)

func TestExtractCity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Asia/Dubai", "Dubai"},
		{"America/New_York", "New York"},
		{"Europe/London", "London"},
		{"UTC", "UTC"},
		{"", "UTC"},
		{"America/Los_Angeles", "Los Angeles"},
		{"Pacific/Auckland", "Auckland"},
		{"Australia/Sydney", "Sydney"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExtractCity(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractCity(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatTimezoneOffset(t *testing.T) {
	tests := []struct {
		offsetSeconds int
		expected      string
	}{
		{14400, "+04:00"},
		{-18000, "-05:00"},
		{0, "+00:00"},
		{3600, "+01:00"},
		{-25200, "-07:00"},
		{19800, "+05:30"},
		{20700, "+05:45"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatTimezoneOffset(tt.offsetSeconds)
			if result != tt.expected {
				t.Errorf("FormatTimezoneOffset(%d) = %q, expected %q", tt.offsetSeconds, result, tt.expected)
			}
		})
	}
}

func TestExtractTimezoneFromSymlink(t *testing.T) {
	tests := []struct {
		link     string
		expected string
	}{
		{"/usr/share/zoneinfo/Asia/Dubai", "Asia/Dubai"},
		{"../usr/share/zoneinfo/America/New_York", "America/New_York"},
		{"/var/db/timezone/zoneinfo/Europe/London", "Europe/London"},
		{"/usr/share/zoneinfo/UTC", "UTC"},
		{"/no-zoneinfo-here", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.link, func(t *testing.T) {
			result := extractTimezoneFromSymlink(tt.link)
			if result != tt.expected {
				t.Errorf("extractTimezoneFromSymlink(%q) = %q, expected %q", tt.link, result, tt.expected)
			}
		})
	}
}
