package huawei

import "time"

const (
	DefaultBaseURL = "https://intl.fusionsolar.huawei.com"
	PlatformName   = "FusionSolar"
)

const (
	DefaultTimeOut    = 30 * time.Second
	TokenTTL          = 30 * time.Minute
	TokenRefreshAhead = 5 * time.Minute
)

const (
	PathLogin              = "/thirdData/login"              // 登录，获取Token
	PathLogout             = "/thirdData/logout"             // 注销，注销Token
	PathStations           = "/thirdData/stations"           // 获取电站列表
	PathGetDevList         = "/thirdData/getDevList"         // 获取设备列表
	PathGetStationRealKpi  = "/thirdData/getStationRealKpi"  // 获取电站实时数据
	PathGetDevRealKpi      = "/thirdData/getDevRealKpi"      // 获取设备实时数据
	PathGetDevKpiDay       = "/thirdData/getDevKpiDay"       // 获取设备日级别数据
	PathGetDevKpiMonth     = "/thirdData/getDevKpiMonth"     // 获取设备月级别数据
	PathGetDevKpiYear      = "/thirdData/getDevKpiYear"      // 获取设备年级别数据
	PathGetKpiStationHour  = "/thirdData/getKpiStationHour"  // 获取电站小时级数据
	PathGetKpiStationDay   = "/thirdData/getKpiStationDay"   // 获取电站日数据
	PathGetKpiStationMonth = "/thirdData/getKpiStationMonth" // 获取电站月数据
	PathGetKpiStationYear  = "/thirdData/getKpiStationYear"  // 获取电站年数据
	PathGetAlarmList       = "/thirdData/getAlarmList"       // 获取电站活动告警数据
)

const MaxBatch = 100

type FailCode int

const (
	FailCodeOK                   FailCode = 0
	FailCodeNotLogin             FailCode = 305   // 当前不在登录状态
	FailCodeNoPermission         FailCode = 401   // 无数据接口权限
	FailCodeRateLimited          FailCode = 407   // 访问频率过高
	FailCodeSystemBusy           FailCode = 429   // 系统级限流，需间隔1分钟以上进行重试
	FailCodeServerError          FailCode = 20004 // 服务器异常
	FailCodeStationLimitExceeded FailCode = 20015 // 一次最多支持查询100个电站
	FailCodeAccountDisabled      FailCode = 20002 // API账户已禁用
	FailCodeAccountExpired       FailCode = 20003 // API账户过期
	FailCodeLoginFailed          FailCode = 20400 // 登录失败：账号或密码错误 / 参数名错误 / 账户锁定 /
	FailCodeDailyQuotaExceeded   FailCode = 20618 // 单 API 账户单日调用超限
	FailCodeServerBusy           FailCode = 20200 // 系统正忙，请稍后再试
	FailCodeTimeInvalid          FailCode = 20604 // 输入的时间参数有误，开始时间不能大于等于结束时间
	FailCodeTimeNegative         FailCode = 20605 // 时间参数中存在负值
	FailCodeStartTimeFuture      FailCode = 20606 // 只输入开始时间时，开始时间不能大于等于当前时间
)
