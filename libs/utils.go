package libs

import "time"

func GetCurrentMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
