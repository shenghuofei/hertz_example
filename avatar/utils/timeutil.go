package utils

import (
	"fmt"
	"time"
)

const (
	// 常用时间格式
	LayoutDate     = "2006-01-02"
	LayoutDateTime = "2006-01-02 15:04:05"
	LayoutTime     = "15:04:05"
)

// ===================== 基础时间操作 =====================

// Now 返回当前时间（本地时区）
func Now() time.Time {
	return time.Now()
}

// NowUTC 返回当前UTC时间
func NowUTC() time.Time {
	return time.Now().UTC()
}

// GetBeforeDays 获取前N天的时间
func GetBeforeDays(n int) time.Time {
	return time.Now().AddDate(0, 0, -n)
}

// GetAfterDays 获取后N天的时间
func GetAfterDays(n int) time.Time {
	return time.Now().AddDate(0, 0, n)
}

// GetBeforeMonths 获取前N个月的时间
func GetBeforeMonths(n int) time.Time {
	return time.Now().AddDate(0, -n, 0)
}

// GetAfterMonths 获取后N个月的时间
func GetAfterMonths(n int) time.Time {
	return time.Now().AddDate(0, n, 0)
}

// GetBeforeHours 获取前N小时的时间
func GetBeforeHours(n int) time.Time {
	return time.Now().Add(-time.Duration(n) * time.Hour)
}

// GetAfterHours 获取后N小时的时间
func GetAfterHours(n int) time.Time {
	return time.Now().Add(time.Duration(n) * time.Hour)
}

// ===================== 格式化与转换 =====================

// FormatTime 格式化时间为字符串（默认 LayoutDateTime）
func FormatTime(t time.Time, layout ...string) string {
	if len(layout) > 0 {
		return t.Format(layout[0])
	}
	return t.Format(LayoutDateTime)
}

// ParseTime 解析字符串为 time.Time（自动识别格式）
func ParseTime(str string) (time.Time, error) {
	layouts := []string{
		LayoutDateTime,
		LayoutDate,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", str)
}

// ===================== 时间戳转换 =====================

// ToTimestampSeconds 转换为秒级时间戳
func ToTimestampSeconds(t time.Time) int64 {
	return t.Unix()
}

// ToTimestampMilli 转换为毫秒级时间戳
func ToTimestampMilli(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// FromTimestampSeconds 秒级时间戳转 time.Time
func FromTimestampSeconds(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// FromTimestampMilli 毫秒级时间戳转 time.Time
func FromTimestampMilli(ts int64) time.Time {
	return time.Unix(0, ts*1e6)
}

// ===================== 字符串转时间戳 =====================

// StringToTimestampSeconds 将字符串时间转为秒级时间戳
func StringToTimestampSeconds(str string) (int64, error) {
	t, err := ParseTime(str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// StringToTimestampMilli 将字符串时间转为毫秒级时间戳
func StringToTimestampMilli(str string) (int64, error) {
	t, err := ParseTime(str)
	if err != nil {
		return 0, err
	}
	return t.UnixNano() / 1e6, nil
}
