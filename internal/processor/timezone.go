package processor

import (
	"os"
	"strings"
)

func GetTimezone() (iana string, found bool) {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz, true
	}

	if data, err := os.ReadFile("/etc/timezone"); err == nil {
		tz := strings.TrimSpace(string(data))
		if tz != "" {
			return tz, true
		}
	}

	if link, err := os.Readlink("/etc/localtime"); err == nil {
		tz := extractTimezoneFromSymlink(link)
		if tz != "" {
			return tz, true
		}
	}

	return "", false
}

func extractTimezoneFromSymlink(link string) string {
	const zoneinfoSuffix = "/zoneinfo/"
	idx := strings.Index(link, zoneinfoSuffix)
	if idx == -1 {
		return ""
	}
	return link[idx+len(zoneinfoSuffix):]
}

func ExtractCity(iana string) string {
	if iana == "UTC" || iana == "" {
		return "UTC"
	}

	parts := strings.Split(iana, "/")
	if len(parts) < 2 {
		return iana
	}

	city := parts[len(parts)-1]
	city = strings.ReplaceAll(city, "_", " ")
	return city
}

func FormatTimezoneOffset(offsetSeconds int) string {
	sign := "+"
	if offsetSeconds < 0 {
		sign = "-"
		offsetSeconds = -offsetSeconds
	}
	hours := offsetSeconds / 3600
	mins := (offsetSeconds % 3600) / 60
	return sign + padZero(hours) + ":" + padZero(mins)
}

func padZero(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
