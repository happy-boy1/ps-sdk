package service

import "ps-sdk/model"

// Adapter 平台适配器。四个平台的 SDK 差异（鉴权方式、分页字段、设备类型编码）
// 全部收敛在实现内部，上层只依赖本接口。
type Adapter interface {
	// Platform 返回平台元信息
	Platform() Platform
	// ListPowerStations 拉取账号下全部电站，站内 ID 尚未分配（StationID 为 0）
	ListPowerStations() ([]model.PowerStation, error)
	// ListPowerDevices 拉取指定电站下的设备，stationID 为 model.PowerStation.StationID。
	// 返回 *PartialError 时表示只拿到部分设备族，调用方应写入 Devices 并记录 Warnings
	ListPowerDevices(stationID uint64) ([]model.PowerDevice, error)
}

// Credential 平台接入凭据。字段按平台取用：
//
//	solarman  AppID/AppSecret/Account/CountryCode/Password/OrgID
//	sungrow   AppID/AppSecret/Account/Password
//	ginlong   AppID/AppSecret
//	huawei    Account/Password
type Credential struct {
	PlatformID int16
	Code       string

	APIURL    string // 接口地址，空则用平台默认地址
	AppID     string // 应用 Key / APIID / 华为 API 账户名
	AppSecret string // 应用密钥 / APISecret

	Account     string // 登录账号：手机号 / 邮箱 / 用户名
	CountryCode string // 手机号国家码，SolarMan 使用
	Password    string // 登录密码
	OrgID       int64  // 商家 ID，SolarMan 使用
}

// adapterBuilder 平台适配器构造函数
type adapterBuilder func(Credential, Options) (Adapter, error)

// adapterBuilders 已注册的平台构造器，key 为平台短码。
// 新增平台时实现 Adapter 并在此登记即可，上层代码无需改动。
var adapterBuilders = map[string]adapterBuilder{}

// registerAdapter 登记平台适配器构造器
func registerAdapter(code string, builder adapterBuilder) {
	adapterBuilders[normalizeCode(code)] = builder
}

// SupportedPlatforms 返回已实现适配器的平台短码
func SupportedPlatforms() []string {
	out := make([]string, 0, len(adapterBuilders))
	for code := range adapterBuilders {
		out = append(out, code)
	}
	return out
}

// NewAdapter 按凭据构造平台适配器
func NewAdapter(cred Credential, opts Options) (Adapter, error) {
	code := normalizeCode(cred.Code)
	builder, ok := adapterBuilders[code]
	if !ok {
		return nil, wrapErr(cred.PlatformID, cred.Code, "构造适配器", ErrPlatformUnknown)
	}
	return builder(cred, opts.defaults())
}
