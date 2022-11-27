package date

import (
	"time"
)

func CurrentUTCTime() *time.Time {
	currentTime := time.Now()
	time := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), time.UTC)

	return &time
}
