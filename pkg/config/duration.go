package config

import (
	"fmt"
	"strings"
	"time"
)

// Duration 包装 time.Duration，支持 TOML 中的 "30s"、"2m" 这类字符串写法，
// 也兼容直接写纳秒数值。零值按未配置处理，由 Config.normalize 补默认值。
type Duration time.Duration

// UnmarshalText 实现 encoding.TextUnmarshaler，go-toml 会用它解析字符串
func (d *Duration) UnmarshalText(text []byte) error {
	value := strings.TrimSpace(string(text))
	if value == "" {
		*d = 0
		return nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("非法时长 %q: %w", value, err)
	}
	*d = Duration(parsed)
	return nil
}

// MarshalText 序列化回可读的时长字符串
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.Duration().String()), nil
}

// Duration 转回 time.Duration
func (d Duration) Duration() time.Duration { return time.Duration(d) }

// String 返回时长字符串
func (d Duration) String() string { return d.Duration().String() }
