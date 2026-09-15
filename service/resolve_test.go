package service

import (
	"strings"
	"testing"
)

func TestCredentialNormalize(t *testing.T) {
	// 接口地址缺省用平台默认值
	cred := Credential{Code: CodeGinlong}
	cred.normalize()
	if cred.APIURL != "https://api.ginlong.com:13333" {
		t.Fatalf("APIURL = %q", cred.APIURL)
	}

	// SolarMan 缺省国家码，便于手机号登录
	cred = Credential{Code: CodeSolarman}
	cred.normalize()
	if cred.CountryCode != DefaultCountryCode {
		t.Fatalf("CountryCode = %q", cred.CountryCode)
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

func TestCredentialFromEnv(t *testing.T) {
	t.Setenv("PS_SOLARMAN_APPID", "app-id")
	t.Setenv("PS_SOLARMAN_APPSECRET", "app-secret")
	t.Setenv("PS_SOLARMAN_ACCOUNT", "13000000000")
	t.Setenv("PS_SOLARMAN_PASSWORD", "pwd")
	t.Setenv("PS_SOLARMAN_ORGID", "1574")
	t.Setenv("PS_SOLARMAN_APIURL", "https://example.com/")

	cred, err := CredentialFromEnv("SolarMan")
	if err != nil {
		t.Fatalf("CredentialFromEnv: %v", err)
	}
	if cred.PlatformID != PlatformSolarman || cred.Code != CodeSolarman {
		t.Fatalf("平台信息异常: %+v", cred)
	}
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

	if _, err := CredentialFromEnv("not-a-platform"); err == nil {
		t.Fatal("未知平台应返回错误")
	}
}

func TestEnvName(t *testing.T) {
	if got := envName(CodeFusionSolar, "appid"); got != "PS_HUAWEI_APPID" {
		t.Fatalf("envName = %q", got)
	}
	if got := envName("Mixed", "Password"); !strings.HasPrefix(got, "PS_MIXED_") {
		t.Fatalf("envName = %q", got)
	}
}
