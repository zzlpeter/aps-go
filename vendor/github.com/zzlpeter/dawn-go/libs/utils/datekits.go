package utils

import (
	"time"
)

func StampSecond() int64 {
	return time.Now().Unix()
}

func StampNanoSecond() int64 {
	return time.Now().UnixNano()
}

func StampMillSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func TimestampFormat(timestamp int64, format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	timeObj := time.Unix(timestamp, 0)
	return timeObj.Format(format)
}