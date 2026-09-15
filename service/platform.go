package service

import (
	"strings"

	"ps-sdk/model"
)

// 平台 ID。数值与 model.PlatformInfo.PlatformID 一一对应，不要随意调整。
const (
	PlatformSolarman    int16 = 1
	PlatformSungrow     int16 = 2
	PlatformGinlong     int16 = 3
	PlatformFusionSolar int16 = 4
)

// 平台英文短码，与 model.PlatformInfo.PlatformCode 一致，用于按字符串定位平台。
const (
	CodeSolarman    = "solarman"
	CodeSungrow     = "sungrow"
	CodeGinlong     = "ginlong"
	CodeFusionSolar = "huawei"
)

// DefaultCountryCode 手机号登录的默认国家码，目前仅 SolarMan 使用
const DefaultCountryCode = "86"

// Platform 平台元信息
type Platform struct {
	ID         int16  // 平台 ID
	Code       string // 英文短码
	NameEN     string // 英文名称
	NameZH     string // 中文名称
	DefaultURL string // 默认接口地址
}

// Platforms 全部已接入平台，顺序即平台 ID 顺序
var Platforms = []Platform{
	{ID: PlatformSolarman, Code: CodeSolarman, NameEN: "solarman", NameZH: "小麦智电", DefaultURL: "https://api.solarmanpv.com"},
	{ID: PlatformSungrow, Code: CodeSungrow, NameEN: "sungrow", NameZH: "阳光云", DefaultURL: "https://gateway.isolarcloud.com"},
	{ID: PlatformGinlong, Code: CodeGinlong, NameEN: "solis", NameZH: "锦浪云", DefaultURL: "https://api.ginlong.com:13333"},
	{ID: PlatformFusionSolar, Code: CodeFusionSolar, NameEN: "fusionsolar", NameZH: "华为 FusionSolar", DefaultURL: "https://intl.fusionsolar.huawei.com"},
}

// PlatformByCode 按英文短码查平台，忽略大小写与分隔符
func PlatformByCode(code string) (Platform, bool) {
	for _, p := range Platforms {
		if p.Code == normalizeCode(code) {
			return p, true
		}
	}
	return Platform{}, false
}

// PlatformByID 按平台 ID 查平台
func PlatformByID(id int16) (Platform, bool) {
	for _, p := range Platforms {
		if p.ID == id {
			return p, true
		}
	}
	return Platform{}, false
}

// normalizeCode 去掉分隔符并转小写，便于容忍 "FusionSolar"、"fusion-solar" 等写法
func normalizeCode(code string) string {
	replacer := strings.NewReplacer("-", "", "_", "", " ", "")
	return strings.ToLower(replacer.Replace(code))
}

// PlatformInfos 返回可直接写入 platform_info 表的记录
func PlatformInfos() []model.PlatformInfo {
	infos := make([]model.PlatformInfo, 0, len(Platforms))
	for _, p := range Platforms {
		infos = append(infos, model.PlatformInfo{
			PlatformID:     p.ID,
			PlatformCode:   p.Code,
			PlatformNameEn: p.NameEN,
			PlatformNameZh: p.NameZH,
		})
	}
	return infos
}
