package service

import (
	"testing"
)

func TestCredentialNormalize(t *testing.T) {
	// 接口地址缺省用平台默认值
	cred := Credential{Code: CodeGinlong}
	cred.normalize()
	if cred.APIURL != "https://api.ginlong.com:13333" {
		t.Fatalf("APIURL = %q", cred.APIURL)
	}

	// SolarMan 以手机号登录时补齐国家码
	cred = Credential{Code: CodeSolarman, Account: "13000000000"}
	cred.normalize()
	if cred.CountryCode != DefaultCountryCode {
		t.Fatalf("CountryCode = %q", cred.CountryCode)
	}

	// 用户名登录不应注入国家码
	cred = Credential{Code: CodeSolarman, Account: "someone"}
	cred.normalize()
	if cred.CountryCode != "" {
		t.Fatalf("用户名登录不应有 CountryCode: %q", cred.CountryCode)
	}

	// 华为把 app_key/app_secret 复用为账户名与密码
	cred = Credential{Code: CodeFusionSolar, AppID: "apiUser", AppSecret: "pwd"}
	cred.normalize()
	if cred.Account != "apiUser" || cred.Password != "pwd" {
		t.Fatalf("华为凭据未派生: %+v", cred)
	}

	// 显式填写的账户名与密码不被覆盖
	cred = Credential{Code: CodeFusionSolar, AppID: "key", AppSecret: "secret", Account: "user", Password: "pass"}
	cred.normalize()
	if cred.Account != "user" || cred.Password != "pass" {
		t.Fatalf("显式凭据被覆盖: %+v", cred)
	}
}

func TestCredentialEnvOverride(t *testing.T) {
	t.Setenv("PS_SOLARMAN_APPID", "app-id")
	t.Setenv("PS_SOLARMAN_APPSECRET", "app-secret")
	t.Setenv("PS_SOLARMAN_ACCOUNT", "13000000000")
	t.Setenv("PS_SOLARMAN_PASSWORD", "pwd")
	t.Setenv("PS_SOLARMAN_ORGID", "1574")
	t.Setenv("PS_SOLARMAN_APIURL", "https://example.com/")

	cred := Credential{PlatformID: PlatformSolarman, Code: CodeSolarman}
	cred.applyEnv()

	if cred.AppID != "app-id" || cred.AppSecret != "app-secret" {
		t.Fatalf("AppID/AppSecret 异常: %+v", cred)
	}
	if cred.Account != "13000000000" || cred.Password != "pwd" || cred.OrgID != 1574 {
		t.Fatalf("登录凭据异常: %+v", cred)
	}
	if cred.APIURL != "https://example.com/" {
		t.Fatalf("APIURL 未被环境变量覆盖: %q", cred.APIURL)
	}
	if cred.CountryCode != DefaultCountryCode {
		t.Fatalf("CountryCode = %q", cred.CountryCode)
	}
	if cred.Identity() != "13000000000" {
		t.Fatalf("Identity = %q", cred.Identity())
	}
}

// providerStub 模拟 pkg/config 提供的配置来源
type providerStub struct {
	cred Credential
	opts Options
}

func (p providerStub) Credential(code string) (Credential, bool) {
	if p.cred.Code != code {
		return Credential{}, false
	}
	return p.cred, true
}

func (p providerStub) Options() Options { return p.opts }

func TestCredentialFromProvider(t *testing.T) {
	provider = providerStub{cred: Credential{
		Code:        CodeSolarman,
		APIURL:      "https://cfg.example.com",
		Account:     "13000000000",
		Password:    "from-config",
		CountryCode: "86",
	}}
	t.Cleanup(func() { provider = nil })

	cred, err := CredentialFromConfig("solarman")
	if err != nil {
		t.Fatalf("CredentialFromConfig: %v", err)
	}
	if cred.APIURL != "https://cfg.example.com" || cred.Password != "from-config" {
		t.Fatalf("配置未生效: %+v", cred)
	}
	if cred.PlatformID != PlatformSolarman || cred.Code != CodeSolarman {
		t.Fatalf("平台信息异常: %+v", cred)
	}

	// 环境变量优先级高于配置文件
	t.Setenv("PS_SOLARMAN_PASSWORD", "from-env")
	cred, err = CredentialFromConfig("solarman")
	if err != nil {
		t.Fatalf("CredentialFromConfig: %v", err)
	}
	if cred.Password != "from-env" {
		t.Fatalf("环境变量未覆盖配置: %+v", cred)
	}

	if _, err := CredentialFromConfig("not-a-platform"); err == nil {
		t.Fatal("未知平台应返回错误")
	}
}

func TestOptionsFrom(t *testing.T) {
	if got := OptionsFrom(nil); got.Timeout <= 0 || got.Concurrency < 1 {
		t.Fatalf("默认参数未补齐: %+v", got)
	}

	want := Options{Timeout: 5 * 1000 * 1000 * 1000, Concurrency: 4}
	if got := OptionsFrom(providerStub{opts: want}); got.Concurrency != 4 {
		t.Fatalf("配置参数未生效: %+v", got)
	}
}
