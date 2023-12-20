package global

import (
    "time"
)

// CurrentTime returns a consistent current time for the entire application
func CurrentTime() time.Time {
    return time.Now().UTC()
}
