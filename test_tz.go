package main
import (
    "fmt"
    "time"
)

func padZero(n int) string {
if n < 10 {
 "0" + string(rune('0'+n))
}
return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func FormatTimezoneOffset(offsetSeconds int) string {
sign := "+"
if offsetSeconds < 0 {
 = "-"
ds = -offsetSeconds
}
hours := offsetSeconds / 3600
mins := (offsetSeconds % 3600) / 60
return sign + padZero(hours) + ":" + padZero(mins)
}

func main() {
    t := time.Now().Local()
    _, offsetSeconds := t.Zone()
    dateStr := t.Format("Jan 2, 15:04")

    // simulate !found on Windows
    fmt.Printf("Simulated Windows output: %s UTC\n", dateStr)
    fmt.Printf("Actual offset: %s\n", FormatTimezoneOffset(offsetSeconds))
}
