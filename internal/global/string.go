package global

import (
    "slices"
    "strconv"
    "strings"
    "time"

    "github.com/xhit/go-str2duration/v2"
)

// aToB converts a string to a boolean
func aToB(alpha string) bool {
    lower := strings.ToLower(alpha)
    positive := []string{
        "yes",
        "true",
        "on",
    }

    return slices.Contains(positive, lower)
}

// aToD converts a time string to a duration
func aToD(alpha string) time.Duration {
    result, _ := str2duration.ParseDuration(alpha)

    return result
}

// aToI converts a string to an integer
func aToI(alpha string) int {
    result, _ := strconv.Atoi(alpha)

    return result
}

// aToP converts a string to a percentage, e.g. 100 = 1, 15 = 0.15
func aToP(alpha string) float32 {
    return float32(aToI(alpha)) / 100
}
