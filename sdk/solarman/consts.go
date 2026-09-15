package solarman

import "time"

// 平台常量
const (
	PlatformName = "Solarman"
	// BaseURLChina 中国数据中心
	BaseURLChina = "https://api.solarmanpv.com"
	// BaseURLGlobal 国际数据中心
	BaseURLGlobal = "https://globalapi.solarmanpv.com"
	// DefaultBaseURL 默认数据中心
	DefaultBaseURL = BaseURLChina
)

// 客户端默认参数
const (
	DefaultTimeout     = 30 * time.Second
	DefaultMaxAttempts = 2
	TokenRefreshAhead  = 5 * time.Minute // 提前重新获取 Token 的窗口
	TokenDefaultTTL    = 60 * 24 * time.Hour
)

// 请求头与查询参数
const (
	HeaderAuthorization = "Authorization" // 值形如 "bearer {access_token}"
	BearerPrefix        = "bearer "       // 注意小写且带一个空格
	QueryAppID          = "appId"
	QueryAppSecret      = "appSecret" // 部分接口（2.4/2.8/5.1）用 query 传密钥鉴权
	QueryLanguage       = "language"
)

// 语言
const (
	LangZh = "zh"
	LangEn = "en"
)

// 请求路径
const (
	// 账号接口
	PathToken       = "/account/v1.0/token"            // 2.1 获取 Token
	PathInfo        = "/account/v1.0/info"             // 2.2 账号商家关系
	PathRole        = "/account/v1.0/role"             // 2.3 账号内权限
	PathUser        = "/account/v1.0/user"             // 2.4 注册帐号
	PathUserInfo    = "/account/v1.0/user-info"        // 2.5 获取账号信息
	PathUserInfoUpd = "/account/v1.0/user-info-update" // 2.6 修改账号信息
	PathBindInfo    = "/account/v1.0/bind-info"        // 2.7 修改绑定信息
	PathPasswordRst = "/account/v1.0/password-reset"   // 2.8 重置密码
	PathPasswordUpd = "/account/v1.0/password-update"  // 2.9 修改密码
	PathCancelCheck = "/account/v1.0/cancelCheck"      // 2.10 账号注销校验
	PathCancel      = "/account/v1.0/cancel"           // 2.11 账号注销
	PathCaptcha     = "/account/v1.0/captcha"          // 5.1 生成验证码
	PathBalance     = "/account/v1.0/balance"          // 13 查询 APPID 剩余可调用次数
)

// 设备接口路径
const (
	PathAlertDetail    = "/device/v1.0/alertDetail"   // 3.1 设备报警明细
	PathAlertList      = "/device/v1.0/alertList"     // 3.2 设备报警列表
	PathCurrentData    = "/device/v1.0/currentData"   // 3.3 设备实时数据
	PathHistorical     = "/device/v1.0/historical"    // 3.4 设备历史数据
	PathDeviceList     = "/device/v1.0/list"          // 3.5 设备库-设备列表
	PathSimInfo        = "/device/v1.0/simInfo"       // 3.6 网关设备 SIM 卡信息
	PathCommunication  = "/device/v1.0/communication" // 3.8 设备通讯关系
	PathCustomControl  = "/device/v1.0/customControl" // 3.9 自定义指令（指令透传）
	PathDeviceRegister = "/device/v1.0/register"      // 4.13 添加网关（采集器或 DTU）
	PathDeviceDelete   = "/device/v1.0/delete"        // 4.14 删除设备
	PathDeviceUpgrade  = "/device/v1.0/upgrade"       // 设备固件升级
)

// 电站接口路径
const (
	PathStationBase     = "/station/v1.0/base"         // 4.1 查询电站基础信息
	PathStationDevice   = "/station/v1.0/device"       // 4.2 获取电站设备列表
	PathStationHistory  = "/station/v1.0/history"      // 4.3 获取电站历史数据
	PathStationList     = "/station/v1.0/list"         // 4.4 获取账号下电站列表
	PathStationRealTime = "/station/v1.0/realTime"     // 4.5 获取电站实时数据
	PathStationRole     = "/station/v1.0/role"         // 4.6 获取电站操作权限
	PathStationAlert    = "/station/v1.0/alert"        // 4.7 获取电站报警列表
	PathStationCreate   = "/station/v1.0/create"       // 4.8 创建电站
	PathStationUpdate   = "/station/v1.0/update"       // 4.9 修改电站
	PathStationDelete   = "/station/v1.0/delete"       // 4.10 删除电站
	PathStationMetering = "/station/v1.0/metering"     // 4.11 累计发电量的计算方法设置
	PathStationOffset   = "/station/v1.0/offset"       // 4.12 设置偏移量
	PathStationAlertV2  = "/station/v1.0/alertV2"      // 4.16 电站报警列表 V2.0
	PathStationWeather  = "/station/v1.0/weather/info" // 4.17 查询电站实时天气
)

// 天气接口路径
const (
	PathWeatherCurrent  = "/weather/v1.0/getCurrentWeather"  // 4.18 查询实时天气
	PathWeatherForecast = "/weather/v1.0/getForecastWeather" // 4.19 查询未来天气
)

// Code 响应信息码，成功时为 null
type Code string

// 响应码（文档 6 响应状态值说明）
const (
	CodeOK                  = "1000000" // success
	CodeServiceUnavailable  = "3201001" // Service Currently Unavailable
	CodeTooFrequently       = "2101002" // client execute Too Frequently
	CodeForbidden           = "2101003" // client is forbidden
	CodeRPCException        = "3501004" // remote rpc exception
	CodeParamNotExist       = "2101005" // param not exist
	CodeInvalidParam        = "2101006" // invalid param
	CodeListTooLong         = "2101007" // list too long
	CodeInvalidParamType    = "2101008" // invalid param type
	CodeAppIDLocked         = "2101009" // appId or api is locked
	CodeInsufficientAllow   = "2101010" // appId insufficient allowance
	CodeUsersUpToFive       = "2101011" // number of users up to five
	CodeWithin30Days        = "2101012" // should be within 30 days
	CodeWithin12Month       = "2101013" // should be within 12 month
	CodeStartAfterEnd       = "2101014" // start time should be earlier than end time
	CodeTimeFormatError     = "2101015" // time format error
	CodeNoUploadRecords     = "2101016" // device no upload records found
	CodeTokenNotFound       = "2101017" // auth token not found
	CodeAppIDNotFound       = "2101018" // auth appId not found
	CodeInvalidToken        = "2101019" // auth invalid token
	CodeInvalidOrgID        = "2101020" // auth invalid orgId
	CodeInvalidAppID        = "2101021" // auth invalid appId
	CodeBusinessError       = "2101022" // bussiness error
	CodeAccessDenied        = "2101023" // access Denied
	CodeResourceNotFound    = "2101024" // resource not found
	CodeAuthFailed          = "2101025" // auth failed
	CodePageTooBig          = "2101026" // page too big
	CodeNoOutPermission     = "2101027" // auth no out operation permission
	CodeInvalidUser         = "2101028" // auth invalid user
	CodeNoOutRoleID         = "2101029" // no out roleId
	CodeNoStationDelPerm    = "2101030" // no out external station delete permission
	CodeStationNotFound     = "2101031" // station not found
	CodeDeviceNotFound      = "2101032" // device not found
	CodeEmailRegistered     = "2101033" // email already registed
	CodePhoneRegistered     = "2101034" // phone number already registed
	CodeEmailNotRegistered  = "2101035" // email not registed
	CodePhoneNotRegistered  = "2101036" // phone number not registed
	CodeUserNameRegistered  = "2101037" // username already registed
	CodeResetNotAllowed     = "2101038" // reset password not allowed
	CodeDateMustBeforeToday = "2101039" // this date can only be before current date
	CodeCaptchaTooOften     = "2102001" // captcha request too frequently
	CodeCaptchaWrong        = "2102002" // captcha was wrong
	CodeCaptchaExpired      = "2102003" // captcha was expire
	CodeOldPasswordWrong    = "2102008" // old password wrong
	CodeUserNameFormat      = "2102009" // username format error
	CodeEmailFormat         = "2102010" // email format error
	CodePhoneFormat         = "2102011" // phone number format error
	CodeGatewayExisted      = "2104001" // gateway existed
	CodeAutoGatewayExisted  = "2104002" // auto discovery gateway existed
	CodeGatewayNotAuto      = "2104003" // gateway can not be auto discovery
)

// codeHints 响应码中文提示
var codeHints = map[string]string{
	CodeOK:                  "成功",
	CodeServiceUnavailable:  "服务暂不可用，请稍后重试",
	CodeTooFrequently:       "调用过于频繁（默认限制 300 次/10 秒），请降低频率",
	CodeForbidden:           "客户端被禁止访问",
	CodeRPCException:        "远程 RPC 异常，平台可能正在更新，请稍后重试",
	CodeParamNotExist:       "缺少参数",
	CodeInvalidParam:        "参数无效，检查大小写、必填项与参数位置",
	CodeListTooLong:         "列表过长",
	CodeInvalidParamType:    "参数类型无效",
	CodeAppIDLocked:         "appId 或接口被锁定：修改数据的接口需二次开通权限，请联系商务",
	CodeInsufficientAllow:   "appId 调用次数已用完，请联系商务或 customerservice@solarmanpv.com 充值",
	CodeUsersUpToFive:       "用户数已达上限 5 个",
	CodeWithin30Days:        "时间范围不能超过 30 天",
	CodeWithin12Month:       "时间范围不能超过 12 个月",
	CodeStartAfterEnd:       "开始时间必须早于结束时间",
	CodeTimeFormatError:     "时间格式错误",
	CodeNoUploadRecords:     "设备无上传记录",
	CodeTokenNotFound:       "缺少 Authorization 头，或未加 bearer 前缀",
	CodeAppIDNotFound:       "appId 不存在",
	CodeInvalidToken:        "access_token 无效或已过期，需重新获取",
	CodeInvalidOrgID:        "orgId 无效",
	CodeInvalidAppID:        "appId 无效",
	CodeBusinessError:       "请求越权：账号与使用端（Home/Pro）需保持一致",
	CodeAccessDenied:        "无权访问该数据，检查 token 与账号权限",
	CodeResourceNotFound:    "资源不存在",
	CodeAuthFailed:          "认证失败",
	CodePageTooBig:          "分页过大",
	CodeNoOutPermission:     "无对外操作权限",
	CodeInvalidUser:         "用户无效",
	CodeNoOutRoleID:         "缺少 roleId",
	CodeNoStationDelPerm:    "无删除外部电站的权限",
	CodeStationNotFound:     "电站不存在",
	CodeDeviceNotFound:      "设备不存在",
	CodeEmailRegistered:     "邮箱已被注册",
	CodePhoneRegistered:     "手机号已被注册",
	CodeEmailNotRegistered:  "邮箱未注册",
	CodePhoneNotRegistered:  "手机号未注册",
	CodeUserNameRegistered:  "用户名已被注册",
	CodeResetNotAllowed:     "不允许重置密码",
	CodeDateMustBeforeToday: "该日期只能早于当前日期",
	CodeCaptchaTooOften:     "验证码请求过于频繁",
	CodeCaptchaWrong:        "验证码错误",
	CodeCaptchaExpired:      "验证码已过期",
	CodeOldPasswordWrong:    "原密码错误",
	CodeUserNameFormat:      "用户名格式错误",
	CodeEmailFormat:         "邮箱格式错误",
	CodePhoneFormat:         "手机号格式错误",
	CodeGatewayExisted:      "网关已存在",
	CodeAutoGatewayExisted:  "自动发现的网关已存在",
	CodeGatewayNotAuto:      "该网关不支持自动发现",
}

// CodeHint 响应码提示
func CodeHint(code string) string {
	if h, ok := codeHints[code]; ok {
		return h
	}
	return "详见 SolarMan OpenAPI 文档「6 响应状态值说明」"
}

// GridInterconnectionType 并网类型（枚举 7）
type GridInterconnectionType string

const (
	GridDistributedFully GridInterconnectionType = "DISTRIBUTED_FULLY" // 光伏+电网
	GridExcess           GridInterconnectionType = "EXCESS"            // 光伏+电网+用电
	GridBatteryBackup    GridInterconnectionType = "BATTERY_BACKUP"    // 光伏+电网+用电+储能
)

// StationType 电站类型（枚举 7）
type StationType string

const (
	StationHouseRoof      StationType = "HOUSE_ROOF"      // 分布式户用
	StationCommercialRoof StationType = "COMMERCIAL_ROOF" // 分布式商业
	StationIndustrialRoof StationType = "INDUSTRIAL_ROOF" // 分布式工业
	StationGround         StationType = "GROUND"          // 地面电站
)

// DeviceType 设备类型（枚举 7）
type DeviceType string

const (
	DeviceInverter       DeviceType = "INVERTER"        // 逆变器
	DeviceWeatherStation DeviceType = "WEATHER_STATION" // 气象站
	DeviceMeter          DeviceType = "METER"           // 电表
	DeviceDTU            DeviceType = "DTU"             // DTU
	DeviceFan            DeviceType = "FAN"             // 风机
	DeviceCollector      DeviceType = "COLLECTOR"       // 采集器
	DevicePVModule       DeviceType = "PV_MODULE"       // PV
	DeviceMicroInverter  DeviceType = "MICRO_INVERTER"  // 微逆
	DeviceBattery        DeviceType = "BATTERY"         // 电池
	DeviceRepeater       DeviceType = "REPEATER"        // REPEATER
)

// DeviceState 设备状态
type DeviceState int

const (
	DeviceOnline  DeviceState = 1 // 在线
	DeviceAlarm   DeviceState = 2 // 报警
	DeviceOffline DeviceState = 3 // 离线
)

func (s DeviceState) String() string {
	switch s {
	case DeviceOnline:
		return "在线"
	case DeviceAlarm:
		return "报警"
	case DeviceOffline:
		return "离线"
	default:
		return "未知"
	}
}

// TotalProductionType 累计发电量计算方法
type TotalProductionType int

const (
	TotalByDeviceTotal TotalProductionType = 1 // 通过设备上传的累计数据计算
	TotalByDailySum    TotalProductionType = 2 // 通过设备上传的日数据求和计算
)
