package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
)

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFileName)
	content := `
[http]
timeout = "5s"
max_attempts = 7
retry_delay = "500ms"
debug = true

[sync]
concurrency = 4

[database]
host = "db.example.com"
port = 3307
user = "u"
password = "p"
name = "n"

[platform.solarman]
api_url = "https://solarman.example.com"
account = "13000000000"
password = "secret"
org_id = 1574
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写入临时配置失败: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Path != path {
		t.Fatalf("Path = %q", cfg.Path)
	}
	if cfg.HTTP.Timeout.Duration() != 5*time.Second || cfg.HTTP.MaxAttempts != 7 {
		t.Fatalf("http 段解析异常: %+v", cfg.HTTP)
	}
	if cfg.HTTP.RetryDelay.Duration() != 500*time.Millisecond || !cfg.HTTP.Debug {
		t.Fatalf("http 段解析异常: %+v", cfg.HTTP)
	}
	if cfg.Sync.Concurrency != 4 {
		t.Fatalf("sync 段解析异常: %+v", cfg.Sync)
	}
	if cfg.Database.Host != "db.example.com" || cfg.Database.Port != 3307 {
		t.Fatalf("database 段解析异常: %+v", cfg.Database)
	}
	if got := cfg.Database.DSNOrBuild(); got != "u:p@tcp(db.example.com:3307)/n?charset=utf8mb4&parseTime=True&loc=Local" {
		t.Fatalf("DSN = %q", got)
	}

	secrets, ok := cfg.Secrets("solarman")
	if !ok {
		t.Fatal("未取到 solarman 配置")
	}
	if secrets.APIURL != "https://solarman.example.com" || secrets.Account != "13000000000" {
		t.Fatalf("平台配置异常: %+v", secrets)
	}
	if secrets.OrgID != 1574 || secrets.Password != "secret" {
		t.Fatalf("平台配置异常: %+v", secrets)
	}

	// 服务层接口：凭据与运行参数
	cred, ok := cfg.Credential("SolarMan")
	if !ok {
		t.Fatal("Credential 未命中")
	}
	if cred.Account != "13000000000" || cred.OrgID != 1574 {
		t.Fatalf("服务层凭据异常: %+v", cred)
	}
	opts := cfg.Options()
	if opts.Timeout != 5*time.Second || opts.Concurrency != 4 || !opts.Debug {
		t.Fatalf("服务层参数异常: %+v", opts)
	}
}

func TestEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFileName)
	if err := os.WriteFile(path, []byte("[platform.ginlong]\napi_url = \"https://a\"\n"), 0o644); err != nil {
		t.Fatalf("写入临时配置失败: %v", err)
	}

	t.Setenv("PS_DB_HOST", "env-host")
	t.Setenv("PS_HTTP_TIMEOUT", "9s")
	t.Setenv("PS_SYNC_CONCURRENCY", "8")
	t.Setenv("PS_GINLONG_APIURL", "https://b")
	t.Setenv("PS_DEBUG", "true")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Database.Host != "env-host" {
		t.Fatalf("数据库环境变量未生效: %q", cfg.Database.Host)
	}
	if cfg.HTTP.Timeout.Duration() != 9*time.Second || !cfg.HTTP.Debug {
		t.Fatalf("http 环境变量未生效: %+v", cfg.HTTP)
	}
	if cfg.Sync.Concurrency != 8 {
		t.Fatalf("sync 环境变量未生效: %+v", cfg.Sync)
	}
	if secrets, _ := cfg.Secrets("ginlong"); secrets.APIURL != "https://b" {
		t.Fatalf("平台环境变量未生效: %+v", secrets)
	}
}

func TestLoadMissingFileWritesTemplate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", DefaultFileName)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Path != path {
		t.Fatalf("Path = %q", cfg.Path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("未写出配置模板: %v", err)
	}

	// 模板必须能被解析，且写入后再次加载结果一致
	again, err := Load(path)
	if err != nil {
		t.Fatalf("重新加载模板失败: %v", err)
	}
	if again.Database.Name != cfg.Database.Name || again.HTTP.MaxAttempts != cfg.HTTP.MaxAttempts {
		t.Fatalf("两次加载结果不一致: %+v / %+v", cfg, again)
	}
}

// TestExampleTemplate 保证随代码分发的模板可解析且覆盖全部平台
func TestExampleTemplate(t *testing.T) {
	var cfg Config
	if err := toml.Unmarshal([]byte(Template()), &cfg); err != nil {
		t.Fatalf("配置模板无法解析: %v", err)
	}
	cfg.normalize()

	for _, code := range []string{"solarman", "sungrow", "ginlong", "huawei"} {
		secrets, ok := cfg.Secrets(code)
		if !ok {
			t.Fatalf("模板缺少 [platform.%s]", code)
		}
		if secrets.APIURL == "" {
			t.Fatalf("[platform.%s] 缺少 api_url", code)
		}
	}
	if cfg.HTTP.Timeout <= 0 || cfg.Database.Port == 0 {
		t.Fatalf("模板默认值异常: %+v", cfg)
	}
}

func TestDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultFileName)
	if err := os.WriteFile(path, []byte("# 空配置\n"), 0o644); err != nil {
		t.Fatalf("写入临时配置失败: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTP.Timeout.Duration() != defaultHTTPTimeout || cfg.Sync.Concurrency != 1 {
		t.Fatalf("默认值未补齐: %+v", cfg)
	}
	if cfg.Database.DSNOrBuild() == "" || cfg.Database.ConnMaxLifeDuration() != time.Hour {
		t.Fatalf("数据库默认值异常: %+v", cfg.Database)
	}
	if cfg.Database.MaxOpen != 100 || cfg.Database.MaxIdle != 12 {
		t.Fatalf("连接池默认值异常: %+v", cfg.Database)
	}
}

func TestPlatformGet(t *testing.T) {
	p := Platform{"api_url": " u ", "org_id": "12", "bad": "x"}
	if got := p.Get("API_URL"); got != "u" {
		t.Fatalf("Get 未忽略大小写或未去空格: %q", got)
	}
	if got := p.Int64("org_id"); got != 12 {
		t.Fatalf("Int64 = %d", got)
	}
	if got := p.Int64("bad"); got != 0 {
		t.Fatalf("非法整数应返回 0，实际 %d", got)
	}
	if p.Enabled() {
		t.Fatal("未配置 app_id/app_secret 时 Enabled 应为 false")
	}
	if !(Platform{"app_id": "x"}).Enabled() {
		t.Fatal("配置了 app_id 时 Enabled 应为 true")
	}
}
