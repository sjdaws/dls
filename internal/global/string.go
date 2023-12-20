package global

import (
    "slices"
    "strconv"
    "strings"
)

func aToB(alpha string) bool {
    lower := strings.ToLower(alpha)
    positive := []string{
        "yes",
        "true",
        "on",
    }

    return slices.Contains(positive, lower)
}

func aToI(alpha string) int {
    result, _ := strconv.Atoi(alpha)

    return result
}
