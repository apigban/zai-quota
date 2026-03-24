package processor

// CalculateWarningLevel returns a warning level based on percentage usage
func CalculateWarningLevel(percentage int) string {
	if percentage < 80 {
		return "safe"
	}
	if percentage < 90 {
		return "warning"
	}
	if percentage < 95 {
		return "critical"
	}
	return "emergency"
}
