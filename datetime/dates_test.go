package datetime

import (
	"testing"
	"time"
)

func TestCheckExpireDays(t *testing.T) {
	now := time.Date(2024, 4, 26, 23, 59, 59, 0, time.Local)
	funcName(t, now)
	now = time.Date(2024, 4, 26, 0, 0, 0, 0, time.Local)
	funcName(t, now)

	//已过期第14天
	days := -14
	now = time.Date(2024, 12, 12, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2024, 11, 28, 16, 28, 55, 0, time.Local)
	b := CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第14天 endTime=%v days=%v now=%v now+days=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat), now.AddDate(0, 0, days).Format(DefaultDateTimeFormat))
	}

	//临过期不足7天当天true
	days = 7
	now = time.Date(2024, 11, 12, 0, 0, 0, 0, time.Local)
	endTime = time.Date(2024, 11, 12, 16, 28, 55, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("临过期不足7天当天true endTime=%v days=%v now=%v now+days=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat), now.AddDate(0, 0, days).Format(DefaultDateTimeFormat))
	}
	//临过期不足7天当天false
	days = 7
	now = time.Date(2024, 11, 12, 18, 0, 0, 0, time.Local)
	endTime = time.Date(2024, 11, 12, 16, 28, 55, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if b {
		t.Errorf("临过期不足7天当天false endTime=%v days=%v now=%v now+days=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat), now.AddDate(0, 0, days).Format(DefaultDateTimeFormat))
	}
}

func funcName(t *testing.T, now time.Time) {
	days := 7
	//临过期不足7天
	endTime := time.Date(2024, 5, 3, 23, 59, 59, 0, time.Local)
	b := CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("临过期不足7天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 5, 3, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("临过期不足7天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//临过期不足6天
	endTime = time.Date(2024, 5, 2, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("临过期不足6天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 5, 2, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("临过期不足6天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//临过期不足8天
	endTime = time.Date(2024, 5, 4, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if b {
		t.Errorf("临过期不足8天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 5, 4, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if b {
		t.Errorf("临过期不足8天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}

	days = -1
	//已过期第1天
	endTime = time.Date(2024, 4, 25, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第1天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 25, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第1天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	days = -7
	//已过期第7天
	endTime = time.Date(2024, 4, 19, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第7天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 19, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第7天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//已过期第6天
	endTime = time.Date(2024, 4, 20, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if b {
		t.Errorf("已过期第6天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 20, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if b {
		t.Errorf("已过期第6天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//已过期第8天
	endTime = time.Date(2024, 4, 18, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if b {
		t.Errorf("已过期第8天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 18, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if b {
		t.Errorf("已过期第8天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	days = -14
	//已过期第14天
	endTime = time.Date(2024, 4, 12, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第14天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 12, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, false)
	if !b {
		t.Errorf("已过期第14天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	days = -14
	//已过期超14天
	endTime = time.Date(2024, 4, 12, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("已过期超14天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 12, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("已过期超14天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//已过期超13天
	endTime = time.Date(2024, 4, 13, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if b {
		t.Errorf("已过期超13天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 13, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if b {
		t.Errorf("已过期超13天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	//已过期超15天
	endTime = time.Date(2024, 4, 11, 23, 59, 59, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("已过期超15天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
	endTime = time.Date(2024, 4, 11, 00, 00, 00, 0, time.Local)
	b = CheckExpireDays(endTime, now, days, true)
	if !b {
		t.Errorf("已过期超15天 endTime=%v days=%v now=%v 测试不通过", endTime.Format(DefaultDateTimeFormat), days, now.Format(DefaultDateTimeFormat))
	}
}
