package functions

import "time"

func TimeDurationToSeconds(duration time.Duration) int {
	return int(duration.Seconds())
}

func Seconds(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
