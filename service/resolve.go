package service

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	"ps-sdk/model"
)

// 环境变量。命名规则 PS_<平台短码大写>_<字段>，例如：
//
//	PS_SOLARMAN_APPID / PS_SOLARMAN_APPSECRET / PS_SOLARMAN_ACCOUNT / PS_SOLARMAN_PASSWORD
//	PS_HUAWEI_APPID 存 API 账户名，PS_HUAWEI_PASSWORD 存密码
//
// 环境变量始终覆盖数据库中的同名字段，便于临时切换或补充表中未建模的字段。
const envPrefix = "PS_"

// CredentialFromDB 从 platform_auth 表读取平台凭据，环境变量可覆盖
func CredentialFromDB(platformID int16) (Credential, error) {
	platform, ok := PlatformByID(platformID)
	if !ok {
		return Credential{}, fmt.Errorf("%w: 平台 ID %d", ErrPlatformUnknown, platformID)
	}
	if db == nil {
		return Credential{}, ErrDatabaseNotReady
	}

	var auth model.PlatformAuth
	err := db.Where("platform_id = ?", platformID).First(&auth).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Credential{}, err
	}

	cred := Credential{
		PlatformID: platformID,
		Code:       platform.Code,
		APIURL:     auth.ApiURL,
		AppID:      auth.AppKey,
		AppSecret:  auth.AppSecret,
	}
	cred.normalize()
	cred.applyEnv()
	return cred, nil
}

// CredentialFromEnv 仅从环境变量读取凭据，便于在没有数据库的脚本中使用
func CredentialFromEnv(code string) (Credential, error) {
	platform, ok := PlatformByCode(code)
	if !ok {
		return Credential{}, fmt.Errorf("%w: %s", ErrPlatformUnknown, code)
	}
	cred := Credential{PlatformID: platform.ID, Code: platform.Code}
	cred.normalize()
	cred.applyEnv()
	return cred, nil
}

// envName 拼出平台字段对应的环境变量名，字段名统一转大写
func envName(code, field string) string {
	return envPrefix + strings.ToUpper(code) + "_" + strings.ToUpper(field)
}

// applyEnv 用环境变量覆盖凭据字段，未设置的环境变量不生效
func (c *Credential) applyEnv() {
	for _, f := range []struct {
		field string
		dest  *string
	}{
		{"APIURL", &c.APIURL},
		{"APPID", &c.AppID},
		{"APPSECRET", &c.AppSecret},
		{"ACCOUNT", &c.Account},
		{"COUNTRYCODE", &c.CountryCode},
		{"PASSWORD", &c.Password},
	} {
		if v := strings.TrimSpace(os.Getenv(envName(c.Code, f.field))); v != "" {
			*f.dest = v
		}
	}
	if v := strings.TrimSpace(os.Getenv(envName(c.Code, "ORGID"))); v != "" {
		c.OrgID = parseInt64(v)
	}
	c.normalize()
}

// normalize 补齐默认值与派生字段：
// 接口地址缺省用平台默认值；华为的 API 账户名/密码存放在 app_id/app_secret 中
func (c *Credential) normalize() {
	if strings.TrimSpace(c.APIURL) == "" {
		if platform, ok := PlatformByCode(c.Code); ok {
			c.APIURL = platform.DefaultURL
		}
	}
	if c.Code == CodeSolarman && strings.TrimSpace(c.CountryCode) == "" {
		c.CountryCode = DefaultCountryCode // SolarMan 手机号登录需要国家码
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
