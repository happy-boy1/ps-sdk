package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// applyEnv 用环境变量覆盖配置，未设置的环境变量不生效
func (c *Config) applyEnv() {
	c.applyHTTPEnv()
	c.applySyncEnv()
	c.applyDatabaseEnv()
	c.applyPlatformEnv()
}

func (c *Config) applyHTTPEnv() {
	setDuration(&c.HTTP.Timeout, "HTTP_TIMEOUT")
	setInt(&c.HTTP.MaxAttempts, "MAX_ATTEMPTS")
	setDuration(&c.HTTP.RetryDelay, "RETRY_DELAY")
	setBool(&c.HTTP.Debug, "DEBUG")
	setBool(&c.HTTP.DevMode, "DEV_MODE")
}

func (c *Config) applySyncEnv() {
	setInt(&c.Sync.Concurrency, "SYNC_CONCURRENCY")
}

func (c *Config) applyDatabaseEnv() {
	setString(&c.Database.DSN, "DB_DSN")
	setString(&c.Database.Host, "DB_HOST")
	setInt(&c.Database.Port, "DB_PORT")
	setString(&c.Database.User, "DB_USER")
	setString(&c.Database.Password, "DB_PASSWORD")
	setString(&c.Database.Name, "DB_NAME")
	setString(&c.Database.Charset, "DB_CHARSET")
	setInt(&c.Database.MaxIdle, "DB_MAX_IDLE")
	setInt(&c.Database.MaxOpen, "DB_MAX_OPEN")
	setString(&c.Database.ConnMaxLife, "DB_CONN_MAX_LIFE")
}

// applyPlatformEnv 覆盖各平台配置，合并进配置文件已有的同名条目
func (c *Config) applyPlatformEnv() {
	if c.Platforms == nil {
		c.Platforms = Platforms{}
	}

	for _, code := range platformCodes {
		platform, key := c.platformEntry(code)

		for field, aliases := range platformEnvAliases {
			value, ok := lookupEnvValue(code, aliases)
			if !ok {
				continue
			}
			platform[field] = value
		}
		c.Platforms[key] = platform
	}
}

// platformEntry 定位平台在配置中的条目，返回条目本身与它在 map 中的键
func (c *Config) platformEntry(code string) (Platform, string) {
	if platform, ok := c.Platforms[code]; ok {
		return platform, code
	}
	normalized := normalizeCode(code)
	for k, platform := range c.Platforms {
		if normalizeCode(k) == normalized {
			return platform, k
		}
	}
	return Platform{}, code
}

// platformCodes 支持环境变量覆盖的平台短码
var platformCodes = []string{"solarman", "sungrow", "ginlong", "huawei"}

// platformEnvAliases 平台字段与可接受的环境变量后缀。
// 采用与 service 层一致的无下划线形式（PS_SOLARMAN_APPID），
// 同时兼容更易读的带下划线形式（PS_SOLARMAN_APP_ID）。
var platformEnvAliases = map[string][]string{
	"api_url":      {"APIURL", "API_URL"},
	"app_id":       {"APPID", "APP_ID"},
	"app_secret":   {"APPSECRET", "APP_SECRET"},
	"account":      {"ACCOUNT"},
	"email":        {"EMAIL"},
	"mobile":       {"MOBILE"},
	"user_name":    {"USERNAME", "USER_NAME"},
	"country_code": {"COUNTRYCODE", "COUNTRY_CODE"},
	"password":     {"PASSWORD"},
	"org_id":       {"ORGID", "ORG_ID"},
}

// lookupEnvValue 依次尝试平台字段的各个环境变量后缀
func lookupEnvValue(code string, aliases []string) (string, bool) {
	for _, alias := range aliases {
		if v, ok := envValue(strings.ToUpper(code) + "_" + alias); ok {
			return v, true
		}
	}
	return "", false
}

// envValue 读取 PS_<key> 环境变量，空值视为未设置
func envValue(key string) (string, bool) {
	v, ok := os.LookupEnv(envPrefix + key)
	if !ok || strings.TrimSpace(v) == "" {
		return "", false
	}
	return strings.TrimSpace(v), true
}

func setString(dest *string, key string) {
	if v, ok := envValue(key); ok {
		*dest = v
	}
}

func setInt(dest *int, key string) {
	v, ok := envValue(key)
	if !ok {
		return
	}
	if n, err := strconv.Atoi(v); err == nil {
		*dest = n
	}
}

func setBool(dest *bool, key string) {
	v, ok := envValue(key)
	if !ok {
		return
	}
	if b, err := strconv.ParseBool(v); err == nil {
		*dest = b
	}
}

func setDuration(dest *Duration, key string) {
	v, ok := envValue(key)
	if !ok {
		return
	}
	if d, err := time.ParseDuration(v); err == nil {
		*dest = Duration(d)
	}
}
