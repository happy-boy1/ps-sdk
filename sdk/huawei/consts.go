package huawei

import "time"

const (
	DefaultBaseURL = "https://intl.fusionsolar.huawei.com"
	PlatformName   = "FusionSolar"
)

const (
	DefaultTimeout    = 30 * time.Second // 单次 HTTP 超时
	TokenTTL          = 30 * time.Minute // XSRF-TOKEN 有效期
	TokenRefreshAhead = 5 * time.Minute  // 提前刷新窗口
	MaxBatch          = 100              // 单次请求最大电站/设备数
	MaxAttempts       = 3                // 默认最大尝试次数
	RetryBaseDelay    = 2 * time.Second  // 重试基础退避
)

// 认证请求头，登录后由响应头下发
const HeaderToken = "XSRF-TOKEN"

// 请求路径
const (
	PathLogin  = "/thirdData/login"
	PathLogout = "/thirdData/logout"
)

// 基础类：基础/监控/告警/报表/配置
const (
	PathStations                 = "/thirdData/stations"
	PathGetDevList               = "/thirdData/getDevList"
	PathGetStationRealKpi        = "/thirdData/getStationRealKpi"
	PathGetDevRealKpi            = "/thirdData/getDevRealKpi"
	PathDeviceHistory            = "/rest/openapi/pvms/nbi/v1/device/history"
	PathGetAlarmList             = "/thirdData/getAlarmList"
	PathGetKpiStationHour        = "/thirdData/getKpiStationHour"
	PathGetKpiStationDay         = "/thirdData/getKpiStationDay"
	PathGetKpiStationMonth       = "/thirdData/getKpiStationMonth"
	PathGetKpiStationYear        = "/thirdData/getKpiStationYear"
	PathGetDevKpiDay             = "/thirdData/getDevKpiDay"
	PathGetDevKpiMonth           = "/thirdData/getDevKpiMonth"
	PathGetDevKpiYear            = "/thirdData/getDevKpiYear"
	PathBatteryModeQuery         = "/rest/openapi/pvms/nbi/v1/configuration/battery-mode"
	PathActivePowerControlMode   = "/rest/openapi/pvms/nbi/v1/configuration/active-power-control-mode"
	PathChargeDischargeTask      = "/rest/openapi/pvms/nbi/v2/control/charge-and-discharge/async-task"
	PathChargeDischargeStatus    = "/rest/openapi/pvms/v1/vpp/chargeAndDischargeStatus"
	PathBatteryModeTask          = "/rest/openapi/pvms/nbi/v1/control/battery/mode/async-task"
	PathBatteryModeTaskInfo      = "/rest/openapi/pvms/nbi/v1/control/battery/mode/task-info"
	PathBatteryConfigTask        = "/rest/openapi/pvms/nbi/v1/control/battery/configuration/async-task"
	PathBatteryConfigTaskInfo    = "/rest/openapi/pvms/nbi/v1/control/battery/configuration/task-info"
	PathActivePowerControlTask   = "/rest/openapi/pvms/nbi/v2/control/active-power-control/async-task"
	PathActivePowerControlTaskIf = "/rest/openapi/pvms/nbi/v2/control/active-power-control/task-info"
	PathBatteryDispatchTask      = "/rest/openapi/pvms/nbi/v1/control/battery/battery-dispatch/async-task"
	PathBatteryDispatchTaskInfo  = "/rest/openapi/pvms/nbi/v1/control/battery/battery-dispatch/task-info"
)

// FailCode 北向接口错误码
type FailCode int

const (
	FailCodeOK                   FailCode = 0     // 正常
	FailCodeTaskPartial          FailCode = 1     // 控制类：任务部分成功
	FailCodeTaskFailed           FailCode = 2     // 控制类：任务失败
	FailCodeNotLogin             FailCode = 305   // 当前不在登录状态
	FailCodeNoPermission         FailCode = 401   // 无数据接口权限
	FailCodeRateLimited          FailCode = 407   // 单用户访问频率过高
	FailCodeSystemBusy           FailCode = 429   // 系统级限流，需间隔 1 分钟以上重试
	FailCodeServerError          FailCode = 20004 // 服务器异常
	FailCodeAccountDisabled      FailCode = 20002 // API 账户已禁用
	FailCodeAccountExpired       FailCode = 20003 // API 账户已过期
	FailCodeLoginFailed          FailCode = 20400 // 登录失败
	FailCodeDailyQuotaExceeded   FailCode = 20618 // 单 API 账户单日调用超限
	FailCodeServerBusy           FailCode = 20200 // 系统正忙
	FailCodeTimeInvalid          FailCode = 20604 // 开始时间不能大于等于结束时间
	FailCodeTimeNegative         FailCode = 20605 // 时间参数中存在负值
	FailCodeStartTimeFuture      FailCode = 20606 // 只传开始时间时不能大于等于当前时间
	FailCodeStationLimitExceeded FailCode = 20015 // 一次最多支持查询 100 个电站
)

// DevTypeID 设备类型 ID
type DevTypeID int

const (
	DevTypeInverter        DevTypeID = 1     // 组串式逆变器
	DevTypeDataCollector   DevTypeID = 2     // 数采
	DevTypeTransformer     DevTypeID = 8     // 箱变
	DevTypeEnvMonitor      DevTypeID = 10    // 环境监测仪
	DevTypeCommGateway     DevTypeID = 13    // 通管机
	DevTypeGeneric         DevTypeID = 16    // 通用设备
	DevTypeGateMeter       DevTypeID = 17    // 关口电表
	DevTypePID             DevTypeID = 22    // PID
	DevTypePinlianDC       DevTypeID = 37    // 品联数采
	DevTypeResidentialInv  DevTypeID = 38    // 户用逆变器
	DevTypeBattery         DevTypeID = 39    // 储能
	DevTypeOffGridCtrl     DevTypeID = 40    // 并离网控制器
	DevTypeESS             DevTypeID = 41    // ESS
	DevTypePLC             DevTypeID = 45    // PLC
	DevTypeOptimizer       DevTypeID = 46    // 优化器
	DevTypeResidentialMtr  DevTypeID = 47    // 户用电表
	DevTypeDongle          DevTypeID = 62    // Dongle
	DevTypeDistributedDC   DevTypeID = 63    // 分布式数采
	DevTypeSafetySwitch    DevTypeID = 70    // 安全关断盒
	DevTypeGrid            DevTypeID = 60001 // 市电
	DevTypeGenerator       DevTypeID = 60003 // 发电机
	DevTypePVModuleGroup   DevTypeID = 60043 // 太阳能模块组
	DevTypePVModule        DevTypeID = 60044 // 太阳能模块
	DevTypePowerConverter  DevTypeID = 60092 // 功率转换器
	DevTypeLithiumCluster  DevTypeID = 60014 // 锂电簇
	DevTypeACDistribution  DevTypeID = 60010 // 交流输出配电
	DevTypeSmartAssistant  DevTypeID = 23070 // SmartAssistant
	DevTypePVESSIntegrated DevTypeID = 23093 // 光储一体机
)

var devTypeNames = map[DevTypeID]string{
	DevTypeInverter:        "组串式逆变器",
	DevTypeDataCollector:   "数采",
	DevTypeTransformer:     "箱变",
	DevTypeEnvMonitor:      "环境监测仪",
	DevTypeCommGateway:     "通管机",
	DevTypeGeneric:         "通用设备",
	DevTypeGateMeter:       "关口电表",
	DevTypePID:             "PID",
	DevTypePinlianDC:       "品联数采",
	DevTypeResidentialInv:  "户用逆变器",
	DevTypeBattery:         "储能",
	DevTypeOffGridCtrl:     "并离网控制器",
	DevTypeESS:             "ESS",
	DevTypePLC:             "PLC",
	DevTypeOptimizer:       "优化器",
	DevTypeResidentialMtr:  "户用电表",
	DevTypeDongle:          "Dongle",
	DevTypeDistributedDC:   "分布式数采",
	DevTypeSafetySwitch:    "安全关断盒",
	DevTypeGrid:            "市电",
	DevTypeGenerator:       "发电机",
	DevTypePVModuleGroup:   "太阳能模块组",
	DevTypePVModule:        "太阳能模块",
	DevTypePowerConverter:  "功率转换器",
	DevTypeLithiumCluster:  "锂电簇",
	DevTypeACDistribution:  "交流输出配电",
	DevTypeSmartAssistant:  "SmartAssistant",
	DevTypePVESSIntegrated: "光储一体机",
}

func (t DevTypeID) String() string {
	if n, ok := devTypeNames[t]; ok {
		return n
	}
	return "未知设备类型"
}

// AlarmLevel 告警级别
type AlarmLevel int

const (
	AlarmLevelCritical AlarmLevel = 1 // 严重
	AlarmLevelMajor    AlarmLevel = 2 // 重要
	AlarmLevelMinor    AlarmLevel = 3 // 次要
	AlarmLevelHint     AlarmLevel = 4 // 提示
)

// Language 告警语言
type Language string

const (
	LangZhCN Language = "zh_CN"
	LangEnUS Language = "en_US"
	LangJaJP Language = "ja_JP"
	LangItIT Language = "it_IT"
	LangNlNL Language = "nl_NL"
	LangPtBR Language = "pt_BR"
	LangDeDE Language = "de_DE"
	LangFrFR Language = "fr_FR"
	LangEsES Language = "es_ES"
	LangPlPL Language = "pl_PL"
	LangRoRO Language = "ro_RO"
	LangCsCZ Language = "cs_CZ"
	LangHuHU Language = "hu_HU"
	LangViVN Language = "vi_VN"
	LangTrTR Language = "tr_TR"
	LangUkUA Language = "uk_UA"
	LangZhTW Language = "zh_TW"
	LangKoKR Language = "ko_KR"
)

// TaskStatus 异步任务状态
type TaskStatus string

const (
	TaskRunning TaskStatus = "RUNNING" // 运行中
	TaskSuccess TaskStatus = "SUCCESS" // 成功
	TaskFail    TaskStatus = "FAIL"    // 失败
)

// TaskFailReason 任务失败原因
type TaskFailReason string

const (
	ReasonFailure   TaskFailReason = "FAILURE"   // 下发失败
	ReasonTimeout   TaskFailReason = "TIMEOUT"   // 超时
	ReasonBusy      TaskFailReason = "BUSY"      // 设备忙
	ReasonInvalid   TaskFailReason = "INVALID"   // 非法设备
	ReasonException TaskFailReason = "EXCEPTION" // 异常
)

// 电站健康状态
const (
	HealthDisconnected = 1 // 断连
	HealthFault        = 2 // 故障
	HealthHealthy      = 3 // 健康
)

// BatteryOperationMode 储能工作模式
type BatteryOperationMode string

const (
	ModeMaximumSelfConsumption BatteryOperationMode = "maximumSelfConsumption" // 最大自发自用
	ModeFullyFidToGrid         BatteryOperationMode = "fullyFidToGrid"         // 全额上网
	ModeTOU                    BatteryOperationMode = "TOU"                    // TOU 模式
	ModeThirdPartyDispatch     BatteryOperationMode = "thirdPartyDispatch"     // 三方调度
	ModeOther                  BatteryOperationMode = "other"                  // 其他
)

// PVEnergyPriority 多余 PV 能量优先级
type PVEnergyPriority string

const (
	PriorityFedToGrid PVEnergyPriority = "fedToGridPreference" // 上网优先
	PriorityCharge    PVEnergyPriority = "chargePreference"    // 充电优先
)

// ChargeOrDischarge 充/放电
type ChargeOrDischarge string

const (
	Charge    ChargeOrDischarge = "charge"    // 充电
	Discharge ChargeOrDischarge = "discharge" // 放电
)

// DispatchSwitch 充放电开关
type DispatchSwitch int

const (
	DispatchStop   DispatchSwitch = 0 // 停止强制充放电
	DispatchCharge DispatchSwitch = 1 // 强制充电
	DispatchDischg DispatchSwitch = 2 // 强制放电
)

// DispatchControlType 充放电控制方式
type DispatchControlType int

const (
	ControlBySOC  DispatchControlType = 1 // SOC 控制
	ControlByTime DispatchControlType = 2 // 时间控制
)
