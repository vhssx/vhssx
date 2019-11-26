package utils

import "time"

func GetMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func GetCurrentMilliseconds() int64 {
	return GetMilliseconds(time.Now())
}
