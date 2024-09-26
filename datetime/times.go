package datetime

import "time"

var DefaultDateTime = time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local)

const DefaultDateTimeFormat = "2006-01-02 15:04:05"
const DefaultDateFormat = "2006-01-02"
const DefaultTimeFormat = "15:04:05"
const DefaultMonthFormat = "2006-01"

// 返回当天0分0秒时刻
func TruncateToDateStart(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
}

// 返回前一天的最后一刻的时间
func TruncateToYesterdayEnd(dt time.Time) time.Time {
	date := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
	return date.Add(-1 * time.Nanosecond)
}
