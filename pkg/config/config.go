package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"

	"ps-sdk/service"
)

// DefaultFileName 默认配置文件名
const DefaultFileName = "config.toml"

// exampleFileName 仓库中随代码分发的配置模板
const exampleFileName = "config.example.toml"

// envPrefix 环境变量前缀
const envPrefix = "PS_"

// ErrNoConfigFile 未找到任何配置文件
var ErrNoConfigFile = errors.New("未找到配置文件")

// 代码内置默认值
const (
	defaultHTTPTimeout = 30 * time.Second
	defaultRetryDelay  = 2 * time.Second
	defaultMaxAttempts = 3
	defaultDBHost      = "127.0.0.1"
	defaultDBPort      = 3306
	defaultDBName      = "rtm"
	defaultDBCharset   = "utf8mb4"
	defaultDBMaxIdle   = 12
	defaultDBMaxOpen   = 100
)

//go:embed config.example.toml
var exampleTemplate []byte

// Template 返回随代码分发的配置模板内容
func Template() string { return string(exampleTemplate) }

// Config 运行配置
type Config struct {
	HTTP      HTTP      `toml:"http"`
	Sync      Sync      `toml:"sync"`
	Database  Database  `toml:"database"`
	Platforms Platforms `toml:"platform"`

	// Path 实际加载的配置文件路径，由 Load 填充
	Path string `toml:"-"`
}

// HTTP 请求相关配置
type HTTP struct {
	Timeout     Duration `toml:"timeout"`      // 单次请求超时，如 "30s"
	MaxAttempts int      `toml:"max_attempts"` // 含首调在内的最大尝试次数
	RetryDelay  Duration `toml:"retry_delay"`  // 重试基础退避，如 "2s"
	Debug       bool     `toml:"debug"`        // 打印请求日志
	DevMode     bool     `toml:"dev_mode"`     // 打印完整请求/响应报文
}

// Sync 同步相关配置
type Sync struct {
	Concurrency int `toml:"concurrency"` // 设备同步的并发电站数
}

// Database 数据库连接配置。
// DSN 非空时直接使用，忽略其余分项；否则按 host/port/user/password/name 组装
type Database struct {
	DSN         string `toml:"dsn"` // 形如 user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	Name        string `toml:"name"`
	Charset     string `toml:"charset"`
	MaxIdle     int    `toml:"max_idle_conns"`
	MaxOpen     int    `toml:"max_open_conns"`
	ConnMaxLife string `toml:"conn_max_life"`
}

// Platforms 各平台配置，key 为平台短码（solarman / sungrow / ginlong / huawei）
type Platforms map[string]Platform

// Platform 单个平台的配置项。
// 使用 map 而非固定结构体，便于各平台按需取用自己的字段：
// 平台之间登录身份字段名不统一（account / email / mobile / user_name），
// 由 Credential 负责归一。取值统一走 Get/Int64，自动容忍 TOML 的整数、布尔等类型。
type Platform map[string]any

// Get 读取字符串字段，键名不区分大小写，不存在返回空串
func (p Platform) Get(key string) string {
	if p == nil {
		return ""
	}
	if v, ok := p[key]; ok {
		return stringValue(v)
	}
	for k, v := range p {
		if strings.EqualFold(k, key) {
			return stringValue(v)
		}
	}
	return ""
}

// stringValue 把 TOML 解码出的任意标量转成字符串
func stringValue(v any) string {
	switch value := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(value)
	case bool:
		return strconv.FormatBool(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case int:
		return strconv.Itoa(value)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return strings.TrimSpace(fmt.Sprint(value))
	}
}

// Int64 读取整数字段，非法值返回 0
func (p Platform) Int64(key string) int64 {
	v, err := strconv.ParseInt(p.Get(key), 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// Enabled 平台是否填了 app_id 或 app_secret
func (p Platform) Enabled() bool {
	return p.Get("app_id") != "" || p.Get("app_secret") != ""
}

// PlatformSecrets 平台凭据在配置文件中的形态。
// app_id / app_secret 以数据库 platform_auth 表为准，这里另外提供
// 接口地址与登录身份等数据库未建模的信息。
type PlatformSecrets struct {
	APIURL      string
	AppKey      string
	AppSecret   string
	Account     string
	Email       string
	Mobile      string
	UserName    string
	CountryCode string
	Password    string
	OrgID       int64
}

// Credential 实现 service.ConfigProvider，返回平台在配置文件中登记的凭据
func (c *Config) Credential(code string) (service.Credential, bool) {
	secrets, platform, ok := c.lookup(code)
	if !ok {
		return service.Credential{}, false
	}

	return service.Credential{
		PlatformID:  platform.ID,
		Code:        platform.Code,
		APIURL:      secrets.APIURL,
		AppID:       secrets.AppKey,
		AppSecret:   secrets.AppSecret,
		Account:     secrets.Account,
		Email:       secrets.Email,
		Mobile:      secrets.Mobile,
		UserName:    secrets.UserName,
		CountryCode: secrets.CountryCode,
		Password:    secrets.Password,
		OrgID:       secrets.OrgID,
	}, true
}

// Options 实现 service.ConfigProvider，返回 [http] 与 [sync] 段落的运行参数
func (c *Config) Options() service.Options {
	return service.Options{
		Timeout:     c.HTTP.Timeout.Duration(),
		MaxAttempts: c.HTTP.MaxAttempts,
		RetryDelay:  c.HTTP.RetryDelay.Duration(),
		Debug:       c.HTTP.Debug,
		DevMode:     c.HTTP.DevMode,
		Concurrency: c.Sync.Concurrency,
	}
}

// Secrets 返回平台在配置文件中的凭据，便于调用方自行判断配置是否完整
func (c *Config) Secrets(code string) (PlatformSecrets, bool) {
	secrets, _, ok := c.lookup(code)
	return secrets, ok
}

// lookup 同时返回平台凭据与平台元信息
func (c *Config) lookup(code string) (PlatformSecrets, service.Platform, bool) {
	platform, ok := service.PlatformByCode(code)
	if !ok {
		return PlatformSecrets{}, service.Platform{}, false
	}
	raw, ok := c.Platform(platform.Code)
	if !ok {
		return PlatformSecrets{}, platform, false
	}

	return PlatformSecrets{
		APIURL:      raw.Get("api_url"),
		AppKey:      raw.Get("app_id"),
		AppSecret:   raw.Get("app_secret"),
		Account:     raw.Get("account"),
		Email:       raw.Get("email"),
		Mobile:      raw.Get("mobile"),
		UserName:    raw.Get("user_name"),
		CountryCode: raw.Get("country_code"),
		Password:    raw.Get("password"),
		OrgID:       raw.Int64("org_id"),
	}, platform, true
}

// 编译期断言：*Config 满足服务层要求的配置来源接口
var _ service.ConfigProvider = (*Config)(nil)

// Load 加载配置文件并按环境变量覆盖。
// path 为空时依次尝试 PS_CONFIG、./config.toml、./pkg/config/config.toml；
// 都不存在时写出配置模板并返回默认配置，便于首次运行。
func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	resolved, err := resolvePath(path)
	if err != nil {
		if !errors.Is(err, ErrNoConfigFile) {
			return nil, err
		}
		if err := writeTemplate(resolved); err != nil {
			return nil, err
		}
		cfg.Path = abs(resolved)
		cfg.applyEnv()
		cfg.normalize()
		return cfg, nil
	}

	raw, err := os.ReadFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", resolved, err)
	}
	if err := toml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %w", resolved, err)
	}

	cfg.Path = abs(resolved)
	cfg.applyEnv()
	cfg.normalize()
	return cfg, nil
}

// resolvePath 定位配置文件：
// 显式指定路径或 PS_CONFIG 时以该路径为准（不存在则返回 ErrNoConfigFile），
// 否则依次尝试默认位置，统一返回绝对路径。
func resolvePath(path string) (string, error) {
	for _, explicit := range []string{path, os.Getenv(envPrefix + "CONFIG")} {
		explicit = strings.TrimSpace(explicit)
		if explicit == "" {
			continue
		}
		if _, err := os.Stat(explicit); err != nil {
			return abs(explicit), ErrNoConfigFile
		}
		return abs(explicit), nil
	}

	for _, candidate := range []string{
		filepath.Join(".", DefaultFileName),
		filepath.Join("pkg", "config", DefaultFileName),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return abs(candidate), nil
		}
	}

	// 都不存在时，把模板写到工作目录下的默认位置
	return abs(filepath.Join(".", DefaultFileName)), ErrNoConfigFile
}

// writeTemplate 写出带注释的配置模板，已存在则跳过
func writeTemplate(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("创建配置目录 %s 失败: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, exampleTemplate, 0o644); err != nil {
		return fmt.Errorf("写入配置模板 %s 失败: %w", path, err)
	}
	return nil
}

// abs 转成绝对路径，方便日志定位
func abs(path string) string {
	if absolute, err := filepath.Abs(path); err == nil {
		return absolute
	}
	return path
}

// defaultConfig 返回代码内置的默认配置
func defaultConfig() *Config {
	return &Config{
		HTTP: HTTP{
			Timeout:     Duration(defaultHTTPTimeout),
			MaxAttempts: defaultMaxAttempts,
			RetryDelay:  Duration(defaultRetryDelay),
		},
		Sync: Sync{Concurrency: 1},
		Database: Database{
			Host:        defaultDBHost,
			Port:        defaultDBPort,
			Name:        defaultDBName,
			Charset:     defaultDBCharset,
			MaxIdle:     defaultDBMaxIdle,
			MaxOpen:     defaultDBMaxOpen,
			ConnMaxLife: "1h",
		},
		Platforms: Platforms{},
	}
}

// normalize 补齐零值并收敛非法值
func (c *Config) normalize() {
	if c.HTTP.Timeout <= 0 {
		c.HTTP.Timeout = Duration(defaultHTTPTimeout)
	}
	if c.HTTP.RetryDelay <= 0 {
		c.HTTP.RetryDelay = Duration(defaultRetryDelay)
	}
	if c.HTTP.MaxAttempts < 1 {
		c.HTTP.MaxAttempts = defaultMaxAttempts
	}
	if c.Sync.Concurrency < 1 {
		c.Sync.Concurrency = 1
	}

	db := &c.Database
	if db.Port <= 0 {
		db.Port = defaultDBPort
	}
	if strings.TrimSpace(db.Host) == "" {
		db.Host = defaultDBHost
	}
	if strings.TrimSpace(db.Charset) == "" {
		db.Charset = defaultDBCharset
	}
	if strings.TrimSpace(db.User) == "" {
		db.User = "root"
	}
	if strings.TrimSpace(db.Name) == "" {
		db.Name = defaultDBName
	}
	if db.MaxIdle <= 0 {
		db.MaxIdle = defaultDBMaxIdle
	}
	if db.MaxOpen <= 0 {
		db.MaxOpen = defaultDBMaxOpen
	}
	if c.Platforms == nil {
		c.Platforms = Platforms{}
	}
}

// Platform 按平台短码取配置，键名不区分大小写
func (c *Config) Platform(code string) (Platform, bool) {
	if c == nil || c.Platforms == nil {
		return nil, false
	}
	if p, ok := c.Platforms[code]; ok {
		return p, true
	}
	normalized := normalizeCode(code)
	for k, p := range c.Platforms {
		if normalizeCode(k) == normalized {
			return p, true
		}
	}
	return nil, false
}

// ConnMaxLifeDuration 连接最长存活时间，非法值按 1 小时
func (d Database) ConnMaxLifeDuration() time.Duration {
	v, err := time.ParseDuration(strings.TrimSpace(d.ConnMaxLife))
	if err != nil || v <= 0 {
		return time.Hour
	}
	return v
}

// DSNOrBuild 返回 MySQL DSN：dsn 非空时直接使用，否则按分项组装
func (d Database) DSNOrBuild() string {
	if v := strings.TrimSpace(d.DSN); v != "" {
		return v
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset)
}

// normalizeCode 归一化平台短码，忽略大小写与分隔符
func normalizeCode(code string) string {
	replacer := strings.NewReplacer("-", "", "_", "", " ", "")
	return strings.ToLower(replacer.Replace(strings.TrimSpace(code)))
}
