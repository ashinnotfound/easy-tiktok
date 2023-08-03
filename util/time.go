package util

import "time"

const timeFormat = "2006-01-02 15:04:05"

// TimeFormat //
// 将时间转换为字符串
func TimeFormat(t *time.Time) string {
	return t.Format(timeFormat)
}
