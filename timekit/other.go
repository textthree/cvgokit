package timekit

import (
	"fmt"
	"time"
)

// 验证是否是一个格里高里日期
func Checkdate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}
	return true
}

// 传入字符串日期，加或减去 n 天
// date string 日期格式：2202-03-01
// date int 传正数为加，负数为减去 n 天
func DateStrAddDay(date string, n int) string {
	timestamp, err := Strtotime("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
	}
	timestamp += int64((3600 * 24) * n)
	return TimestampToDate(timestamp)
}

// 判断两个时间是否是相邻的天
func IsAdjacentDays(t1, t2 time.Time) bool {
	// 获取日期部分并设置时间为零点零分零秒
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	t1 = time.Date(y1, m1, d1, 0, 0, 0, 0, t1.Location())
	t2 = time.Date(y2, m2, d2, 0, 0, 0, 0, t2.Location())

	// 计算两个日期之间的差异
	diff := t1.Sub(t2)

	// 判断差异是否为正负一天
	return diff == 24*time.Hour || diff == -24*time.Hour
}

// 睡眠
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

// 延迟执行当前脚本 t 秒
func Usleep(t int64) {
	time.Sleep(time.Duration(t) * time.Microsecond)
}
