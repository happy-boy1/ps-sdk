package service

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	"ps-sdk/model"
)

// 环境变量命名规则 PS_<平台短码大写>_<字段>，例如：
//
//	PS_SOLARMAN_APPID / PS_SOLARMAN_PASSWORD / PS_HUAWEI_APPID / PS_HUAWEI_PASSWORD
//
// 凭据来源优先级由低到高：config.toml → 数据库 platform_auth 表 → 环境变量。
const envPrefix = "PS_"

// envFields 参与环境变量覆盖的凭据字段
var envFields = []struct {
	field string
	dest  func(*Credential) *string
}{
	{"APIURL", func(c *Credential) *string { return &c.APIURL }},
	{"APPID", func(c *Credential) *string { return &c.AppID }},
	{"APPSECRET", func(c *Credential) *string { return &c.AppSecret }},
	{"ACCOUNT", func(c *Credential) *string { return &c.Account }},
	{"EMAIL", func(c *Credential) *string { return &c.Email }},
	{"MOBILE", func(c *Credential) *string { return &c.Mobile }},
	{"USERNAME", func(c *Credential) *string { return &c.UserName }},
	{"COUNTRYCODE", func(c *Credential) *string { return &c.CountryCode }},
	{"PASSWORD", func(c *Credential) *string { return &c.Password }},
}

// CredentialFromDB 组装平台凭据：配置文件提供接口地址与登录身份，
// 数据库 platform_auth 表提供 app_key/app_secret，环境变量可覆盖任意字段。
// 数据库不可用时退化为「配置文件 + 环境变量」。
func CredentialFromDB(platformID int16) (Credential, error) {
	platform, ok := PlatformByID(platformID)
	if !ok {
		return Credential{}, fmt.Errorf("%w: 平台 ID %d", ErrPlatformUnknown, platformID)
	}
	return credential(platform)
}

// CredentialFromConfig 只从配置文件与环境变量组装凭据，不访问数据库
func CredentialFromConfig(code string) (Credential, error) {
	platform, ok := PlatformByCode(code)
	if !ok {
		return Credential{}, fmt.Errorf("%w: %s", ErrPlatformUnknown, code)
	}

	cred := Credential{PlatformID: platform.ID, Code: platform.Code}
	cred.applyProvider()
	cred.normalize()
	cred.applyEnv()
	return cred, nil
}

// credential 按平台组装凭据
func credential(platform Platform) (Credential, error) {
	cred := Credential{PlatformID: platform.ID, Code: platform.Code}
	cred.applyProvider()

	if db != nil {
		var auth model.PlatformAuth
		err := db.Where("platform_id = ?", platform.ID).First(&auth).Error
		switch {
		case err == nil:
			// 数据库是 app_key/app_secret 的权威来源
			cred.AppID = firstNonEmpty(auth.AppKey, cred.AppID)
			cred.AppSecret = firstNonEmpty(auth.AppSecret, cred.AppSecret)
			cred.APIURL = firstNonEmpty(auth.ApiURL, cred.APIURL)
		case errors.Is(err, gorm.ErrRecordNotFound):
			Logf("[%s] platform_auth 无记录，凭据改用配置文件与环境变量", platform.Code)
		default:
			return Credential{}, fmt.Errorf("读取平台凭据失败: %w", err)
		}
	}

	cred.normalize()
	cred.applyEnv()
	return cred, nil
}

// applyProvider 用配置文件中的凭据打底
func (c *Credential) applyProvider() {
	if provider == nil {
		return
	}
	from, ok := provider.Credential(c.Code)
	if !ok {
		return
	}

	c.APIURL = firstNonEmpty(from.APIURL, c.APIURL)
	c.AppID = firstNonEmpty(from.AppID, c.AppID)
	c.AppSecret = firstNonEmpty(from.AppSecret, c.AppSecret)
	c.Account = firstNonEmpty(from.Account, c.Account)
	c.Email = firstNonEmpty(from.Email, c.Email)
	c.Mobile = firstNonEmpty(from.Mobile, c.Mobile)
	c.UserName = firstNonEmpty(from.UserName, c.UserName)
	c.CountryCode = firstNonEmpty(from.CountryCode, c.CountryCode)
	c.Password = firstNonEmpty(from.Password, c.Password)
	if from.OrgID > 0 {
		c.OrgID = from.OrgID
	}
}

// applyEnv 用环境变量覆盖凭据字段，未设置的环境变量不生效
func (c *Credential) applyEnv() {
	for _, f := range envFields {
		if v := envString(c.Code, f.field); v != "" {
			*f.dest(c) = v
		}
	}
	if v := envString(c.Code, "ORGID"); v != "" {
		c.OrgID = parseInt64(v)
	}
	c.normalize()
}

// envString 读取 PS_<平台短码大写>_<字段> 环境变量
func envString(code, field string) string {
	return strings.TrimSpace(os.Getenv(envPrefix + strings.ToUpper(code) + "_" + strings.ToUpper(field)))
}

// normalize 补齐默认值与派生字段：
// 接口地址缺省用平台默认值；SolarMan 手机号登录需要国家码；华为把 API 账户名/密码复用为 app_key/app_secret
func (c *Credential) normalize() {
	if strings.TrimSpace(c.APIURL) == "" {
		if platform, ok := PlatformByCode(c.Code); ok {
			c.APIURL = platform.DefaultURL
		}
	}

	if c.Code == CodeSolarman {
		// 不含 @ 且全为数字时按手机号处理，需要国家码
		if strings.TrimSpace(c.CountryCode) == "" && isDigits(c.Identity()) {
			c.CountryCode = DefaultCountryCode
		}
		return
	}
	if c.Code != CodeFusionSolar {
		return
	}
	if strings.TrimSpace(c.Account) == "" {
		c.Account = strings.TrimSpace(c.AppID)
	}
	if strings.TrimSpace(c.Password) == "" {
		c.Password = strings.TrimSpace(c.AppSecret)
	}
}
