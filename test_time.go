package main
import (
    "fmt"
    "time"
    "math"
)
func main() {
    t := time.Unix(math.MaxInt64, 0).Local()
    d := time.Until(t)
    hours := int(d.Hours())
    days := hours / 24
    fmt.Printf("Days for max int64 seconds: %v\n", days)

    t2 := time.Unix(math.MaxInt32, 0).Local()
    d2 := time.Until(t2)
    fmt.Printf("Days for max int32 seconds (year 2038): %v\n", int(d2.Hours())/24)
    
    // Test parsing zero or -1
    t3 := time.Unix(-1, 0).Local()
    d3 := time.Until(t3)
    fmt.Printf("Days for -1: %v\n", int(d3.Hours())/24)
}
