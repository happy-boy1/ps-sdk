package service

import (
	"strconv"
	"strings"
	"time"
)

// unixSeconds UNIX 秒 → *time.Time，非正值返回 nil
func unixSeconds(sec int64) *time.Time {
	if sec <= 0 {
		return nil
	}
	t := time.Unix(sec, 0)
	return &t
}

// unixMillis UNIX 毫秒 → *time.Time，非正值返回 nil
func unixMillis(ms int64) *time.Time {
	if ms <= 0 {
		return nil
	}
	t := time.UnixMilli(ms)
	return &t
}

// dateTime 按 yyyy-MM-dd 或 RFC3339 解析日期时间，失败返回 nil
func dateTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, value); err == nil {
			return &t
		}
	}
	return nil
}

// parseInt64 宽松解析整数，失败返回 0
func parseInt64(s string) int64 {
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// formatInt64 整数转十进制字符串
func formatInt64(v int64) string { return strconv.FormatInt(v, 10) }

// parseFloat 宽松解析浮点数，失败返回 0
func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return v
}

// float64Ptr 非正数返回 nil，避免把 0 写进带 check 约束的容量字段
func float64Ptr(v float64) *float64 {
	if v <= 0 {
		return nil
	}
	return &v
}

// trimOr 返回去空格后的值，为空时返回 fallback
func trimOr(value, fallback string) string {
	if v := strings.TrimSpace(value); v != "" {
		return v
	}
	return fallback
}

// firstNonEmpty 返回第一个非空值
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// isDigits 判断字符串是否全为数字
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
