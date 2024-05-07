package hello_utils

import "time"

// CheckExpireDays 一般情况判断now+days后年月日相等，daysMore情况增加or小于结束年月日
func CheckExpireDays(endTime, now time.Time, days int, daysMore bool) bool {
	d0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	d := d0.AddDate(0, 0, days)
	t := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())

	//快到期days天多1秒，不能算“不足days天”
	//已过期days天少1秒，不能算“过期days天”
	if d.Before(t) {
		return false
	}

	//不严格校验是否相差days天
	if daysMore {
		if days > 0 {
			//now不能大于endTime，否则就不是临过期，而是已过期
			return d0.Before(t)
		} else {
			//now不能小于endTime，否则就不是已过期，而是临过期
			return d0.After(t)
		}
	}

	//临/已过期第days天，落在同一天即可
	return d.Equal(t)
}
