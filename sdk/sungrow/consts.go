package sungrow

import "time"

const (
	PlatformName   = "Sungrow"
	BaseURLChina   = "https://gateway.isolarcloud.com"
	BaseURLGlobal  = "https://gateway.isolarcloud.com.hk"
	DefaultBaseURL = BaseURLChina
)

const (
	HeaderSysCode    = "sys_code"
	HeaderXAccessKey = "x-access-key"
	QueryAppID       = "appkey"
	QUeryToken       = "token"
	QueryLanguage    = "lang"
)

const (
	MaxAttempts       = 3
	RetryBaseDelay    = 2 * time.Second
	TokenRefreshAhead = 5 * time.Second
	TokenTTL          = 24 * time.Hour
)

const (
	LangZh = "zh"
	LangEn = "en"
)

const (
	PathLogin             = "/openapi/login"
	PathSharePsBySN       = "/openapi/sharePsBySN"
	PathReleasePsSharing  = "/openapi/releasePsSharing"
	PathShareMyPs         = "/openapi/shareMyPs"
	PathCancelPsSharing   = "/openapi/cancelPsSharing"
	PathSharePsBySNAndNmi = "/openapi/sharePowerStationBySnAndNmi"
)

const (
	PathGetPsList                       = "/openapi/getPowerStationList"
	PathGetDevList                      = "/openapi/getDeviceList"
	PathGetPsDetail                     = "/openapi/getPowerStationDetail"
	PathGetDevRtd                       = "/openapi/getDeviceRealTimeData"
	PathGetFaultAlarmInfo               = "/openapi/getFaultAlarmInfo"
	PathGetDevPointMinuteDataList       = "/openapi/getDevicePointMinuteDataList"
	PathGetDevPointDayMonthYearDataList = "/openapi/getDevicePointsDayMonthYearDataList"
	PathGetDevPropertyPointValue        = "/openapi/getDevPropertyPointValue"
	PathGetOpenPointInfo                = "/openapi/getOpenPointInfo"
	PathGetDevListByUser                = "/openapi/getDeviceListByUser"
	PathGetCommunicationDevInfoByDevSN  = "/openapi/getCommunicationDevInfoByDevSn"
	PathGetPvInverterRtd                = "/openapi/getPVInverterRealTimeData"
	PathGetOpenApiCallInfo              = "/openapi/getOpenApiCallInfo"
	PathGetBatchPsDetail                = "/openapi/getBatchPsDetail"
	PathGetDevStringInfo                = "/openapi/getDeviceStringInfo"
	PathGetMlpeDevList                  = "/openapi/getMlpeDeviceList"
	PathGetMlpeRtd                      = "/openapi/getMlpeRealTimeData"
	PathGetMlpeMinuteDataList           = "/openapi/getMlpeMinuteDataList"
	PathGetMlpeDayMonthYearDataList     = "/openapi/getMlpeDayMonthYearDataList"
)

// 暂时不实现
const (
	PathParamSettingCheck          = "/openapi/paramSettingCheck"
	PathParamSetting               = "/openapi/paramSetting"
	PathGetParamSettingTask        = "/openapi/getParamSettingTask"
	PathReadOnlyParamSet           = "/openapi/readOnlyParamSet"
	PathGetReadOnlyResult          = "/openapi/getReadOnlyResult"
	PathGetReadOnlyParamDefinition = "/openapi/getReadOnlyParamDefinition"
	PathGetDevSettingRecordList    = "/openapi/getDeviceSettingRecordList"
)

// 暂时不实现
const (
	PathGetConfig  = "/openapi/datasubscribe/getConfig"
	PathStart      = "/openapi/datasubscribe/start"
	PathStop       = "/openapi/datasubscribe/stop"
	PathGetHisData = "/openapi/datasubscribe/getHisData"
)

type ResultCode string

const (
	ResultCodeSuccess                                                   ResultCode = "1"
	ResultCodeError                                                     ResultCode = "-1"
	ResultCodeErUnknownException                                        ResultCode = "000"
	ResultCodeErMissingParameterAppkey                                  ResultCode = "001"
	ResultCodeErMissingParameterToken                                   ResultCode = "002"
	ResultCodeErMissingParameterSysCode                                 ResultCode = "003"
	ResultCodeErInvalidAppkey                                           ResultCode = "E00000"
	ResultCodeErApiServiceHasExpired                                    ResultCode = "E00001"
	ResultCodeErParameterDecryptError                                   ResultCode = "E00002"
	ResultCodeErTokenLoginInvalid                                       ResultCode = "E00003"
	ResultCodeErMonthCallApiTimesUpperLimit                             ResultCode = "E998"
	ResultCodeErHourCallApiTimesUpperLimit                              ResultCode = "E999"
	ResultCodeErMissingParameter                                        ResultCode = "009"
	ResultCodeErParameterValueInvalid                                   ResultCode = "010"
	ResultCodeErSqlException                                            ResultCode = "011"
	ResultCodeUnauthorizedAccess                                        ResultCode = "E900"
	ResultCodeCallTooFrequently                                         ResultCode = "E901"
	ResultCodeAbnormalNetworkEnvironment                                ResultCode = "E903"
	ResultCodeRequestIsNotEncrypted                                     ResultCode = "E902"
	ResultCodeMissingParameterInRequestHeaderXRandomSecretKey           ResultCode = "E904"
	ResultCodeAesDecryptionException                                    ResultCode = "E905"
	ResultCodeRsaDecryptionException                                    ResultCode = "E906"
	ResultCodeAesRandomSecretKeyLengthMustBe16                          ResultCode = "E907"
	ResultCodeMissingKeyParameterApiKeyParam                            ResultCode = "E908"
	ResultCodeInvalidParameterFormatNonce32BitStringOfNumbersAndLetters ResultCode = "E909"
	ResultCodeRepeatedRequest                                           ResultCode = "E910"
	ResultCodeMissingParameterInRequestHeaderXAccessKey                 ResultCode = "E911"
	ResultCodeIllegalXAccessKey                                         ResultCode = "E912"
	ResultCodeExpiredRequest                                            ResultCode = "E913"
	ResultCodeMismatchedAppkeyAndXAccessKey                             ResultCode = "E914"
	ResultCodeLoginTooFrequently                                        ResultCode = "E916"
	ResultCodePermitDeniedByWhiteIpAddress                              ResultCode = "E918"
	ResultCodePermitDeniedByWhiteListUser                               ResultCode = "E919"
	ResultCodeSystemNotFound                                            ResultCode = "E994"
	ResultCodeRequestBodyTooLarge                                       ResultCode = "E995"
	ResultCodeApiNotFound                                               ResultCode = "E996"
	ResultCodeTransformBusinessResponseDataOccurError                   ResultCode = "E997"
)

var resultCodeHints = map[ResultCode]string{
	ResultCodeSuccess:                                                   "调用成功",
	ResultCodeError:                                                     "服务内部异常",
	ResultCodeErUnknownException:                                        "未知异常",
	ResultCodeErMissingParameterAppkey:                                  "参数appkey不可以为空",
	ResultCodeErMissingParameterToken:                                   "参数token不可以为空",
	ResultCodeErMissingParameterSysCode:                                 "参数sys_code不可以为空",
	ResultCodeErInvalidAppkey:                                           "appkey不合法",
	ResultCodeErApiServiceHasExpired:                                    "API服务到期了",
	ResultCodeErParameterDecryptError:                                   "参数解密异常",
	ResultCodeErTokenLoginInvalid:                                       "token不合法或已失效",
	ResultCodeErMonthCallApiTimesUpperLimit:                             "API月调用次数达上限",
	ResultCodeErHourCallApiTimesUpperLimit:                              "API小时调用次数达上限",
	ResultCodeErMissingParameter:                                        "必传参数缺失",
	ResultCodeErParameterValueInvalid:                                   "参数值不合法",
	ResultCodeErSqlException:                                            "SQL异常",
	ResultCodeUnauthorizedAccess:                                        "未经授权的访问",
	ResultCodeCallTooFrequently:                                         "调用频繁",
	ResultCodeAbnormalNetworkEnvironment:                                "ip地址频繁切换",
	ResultCodeRequestIsNotEncrypted:                                     "请求未加密",
	ResultCodeMissingParameterInRequestHeaderXRandomSecretKey:           "请求头缺失参数：x-random-secret-key",
	ResultCodeAesDecryptionException:                                    "AES解密异常",
	ResultCodeRsaDecryptionException:                                    "RSA解密异常",
	ResultCodeAesRandomSecretKeyLengthMustBe16:                          "AES随机秘钥长度必须为16",
	ResultCodeMissingKeyParameterApiKeyParam:                            "缺失关键参数：api_key_param",
	ResultCodeInvalidParameterFormatNonce32BitStringOfNumbersAndLetters: "参数格式不合法：nonce [数字和字母组合的32位字符串]",
	ResultCodeRepeatedRequest:                                           "重复的请求，请求中的nonce需要重新生成",
	ResultCodeMissingParameterInRequestHeaderXAccessKey:                 "请求头缺失参数：x-access-key",
	ResultCodeIllegalXAccessKey:                                         "非法的 x-access-key",
	ResultCodeExpiredRequest:                                            "请求过期，请求中timestamp（0时区UNIX时间戳）与服务器时间差不在合理范围内",
	ResultCodeMismatchedAppkeyAndXAccessKey:                             "appkey与access-key不匹配",
	ResultCodeLoginTooFrequently:                                        "登录频繁",
	ResultCodePermitDeniedByWhiteIpAddress:                              "白名单 IP 地址拒绝访问",
	ResultCodePermitDeniedByWhiteListUser:                               "白名单用户权限拒绝",
	ResultCodeSystemNotFound:                                            "未找到系统",
	ResultCodeRequestBodyTooLarge:                                       "请求体过大",
	ResultCodeApiNotFound:                                               "api未找到",
	ResultCodeTransformBusinessResponseDataOccurError:                   "转换业务响应数据时发生错误",
}

func (c ResultCode) Hint() string {
	if h, ok := resultCodeHints[c]; ok {
		return h
	}
	return "详细见 Sungrow OpenApi 文档错误码说明"
}

func (c ResultCode) Retryable() bool {
	switch c {
	case ResultCodeRequestIsNotEncrypted, ResultCodeTransformBusinessResponseDataOccurError, ResultCodeSystemNotFound:
		return true
	default:
		return false
	}
}

// 设备类型
type DevTypeID int

const (
	DevTypeInverter                  DevTypeID = 1  // 逆变器
	DevTypeContainer                 DevTypeID = 2  // 集装箱
	DevTypeGridConnectionPoint       DevTypeID = 3  // 并网点
	DevTypeCombinerBox               DevTypeID = 4  // 汇流箱
	DevTypeMeteoStation              DevTypeID = 5  // 环境监测仪
	DevTypeTransformer               DevTypeID = 6  // 变压器
	DevTypeMeter                     DevTypeID = 7  // 电表
	DevTypeUPS                       DevTypeID = 8  // UPS
	DevTypeDataLogger                DevTypeID = 9  // 数据采集器
	DevTypeString                    DevTypeID = 10 // 组串
	DevTypePlant                     DevTypeID = 11 // 电站
	DevTypeCircuitProtection         DevTypeID = 12 // 线路保护
	DevTypeSplittingDevice           DevTypeID = 13 // 解列装置
	DevTypeEnergyStorageSystem       DevTypeID = 14 // 储能逆变器
	DevTypeSamplingDevice            DevTypeID = 15 // 采集设备
	DevTypeEMU                       DevTypeID = 16 // EMU
	DevTypeUnit                      DevTypeID = 17 // 单元
	DevTypeTempHumiditySensor        DevTypeID = 18 // 温湿度传感器
	DevTypeIntelligentPowerCabinet   DevTypeID = 19 // 智能配电柜
	DevTypeDisplayDevice             DevTypeID = 20 // 显示设备
	DevTypeACPowerDistributedCabinet DevTypeID = 21 // 交流配电柜
	DevTypeCommunicationModule       DevTypeID = 22 // 通信模块
	DevTypeSystemBMS                 DevTypeID = 23 // 系统BMS
	DevTypeArrayBMS                  DevTypeID = 24 // 阵列BMS
	DevTypeDCDC                      DevTypeID = 25 // 直流-直流
	DevTypeEnergyManagementSystem    DevTypeID = 26 // 能量管理系统
	DevTypeTrackingSystem            DevTypeID = 27 // 跟踪系统
	DevTypeWindEnergyConverter       DevTypeID = 28 // 风能变流器
	DevTypeSVG                       DevTypeID = 29 // SVG
	DevTypePTCabinet                 DevTypeID = 30 // PT柜
	DevTypeBusProtection             DevTypeID = 31 // 母线保护
	DevTypeCleaningDevice            DevTypeID = 32 // 清扫机器人
	DevTypeDirectCurrentCabinet      DevTypeID = 33 // 直流屏
	DevTypePublicMeasurementControl  DevTypeID = 34 // 公用测控
	DevTypePCS                       DevTypeID = 37 // 储能变流器 (Energiespeichersystem)
	DevTypeOptimizer                 DevTypeID = 41 // 优化器
	DevTypeBattery                   DevTypeID = 43 // 电池
	DevTypeBatteryClusterMgmt        DevTypeID = 44 // 电池簇管理单元
	DevTypeLocalController           DevTypeID = 45 // 本地控制器
	DevTypeCharger                   DevTypeID = 51 // 充电桩
	DevTypeBatterySystemController   DevTypeID = 52 // 电池系统控制器
	DevTypeMicroinverter             DevTypeID = 55 // 微型逆变器
	DevTypeDieselGenerator           DevTypeID = 63 // 柴油发电机
)

var devTypeNames = map[DevTypeID]string{
	DevTypeInverter:                  "逆变器",
	DevTypeContainer:                 "集装箱",
	DevTypeGridConnectionPoint:       "并网点",
	DevTypeCombinerBox:               "汇流箱",
	DevTypeMeteoStation:              "环境监测仪",
	DevTypeTransformer:               "变压器",
	DevTypeMeter:                     "电表",
	DevTypeUPS:                       "UPS",
	DevTypeDataLogger:                "数据采集器",
	DevTypeString:                    "组串",
	DevTypePlant:                     "电站",
	DevTypeCircuitProtection:         "线路保护",
	DevTypeSplittingDevice:           "解列装置",
	DevTypeEnergyStorageSystem:       "储能逆变器",
	DevTypeSamplingDevice:            "采集设备",
	DevTypeEMU:                       "EMU",
	DevTypeUnit:                      "单元",
	DevTypeTempHumiditySensor:        "温湿度传感器",
	DevTypeIntelligentPowerCabinet:   "智能配电柜",
	DevTypeDisplayDevice:             "显示设备",
	DevTypeACPowerDistributedCabinet: "交流配电柜",
	DevTypeCommunicationModule:       "通信模块",
	DevTypeSystemBMS:                 "系统BMS",
	DevTypeArrayBMS:                  "阵列BMS",
	DevTypeDCDC:                      "直流-直流",
	DevTypeEnergyManagementSystem:    "能量管理系统",
	DevTypeTrackingSystem:            "跟踪系统",
	DevTypeWindEnergyConverter:       "风能变流器",
	DevTypeSVG:                       "SVG",
	DevTypePTCabinet:                 "PT柜",
	DevTypeBusProtection:             "母线保护",
	DevTypeCleaningDevice:            "清扫机器人",
	DevTypeDirectCurrentCabinet:      "直流屏",
	DevTypePublicMeasurementControl:  "公用测控",
	DevTypePCS:                       "储能变流器",
	DevTypeOptimizer:                 "优化器",
	DevTypeBattery:                   "电池",
	DevTypeBatteryClusterMgmt:        "电池簇管理单元",
	DevTypeLocalController:           "本地控制器",
	DevTypeCharger:                   "充电桩",
	DevTypeBatterySystemController:   "电池系统控制器",
	DevTypeMicroinverter:             "微型逆变器",
	DevTypeDieselGenerator:           "柴油发电机",
}

func (t DevTypeID) String() string {
	if n, ok := devTypeNames[t]; ok {
		return n
	}
	return "未知类型设备"
}

type PsType uint8

const (
	PsTypeGround             PsType = 1
	PsTypeDistributedPV      PsType = 3
	PsTypeResidentialPV      PsType = 4
	PsTypeResidentialStorage PsType = 5
	PsTypeVillage            PsType = 6
	PsTypeDistributedStorage PsType = 7
	PsTypePovertyRelief      PsType = 8
	PsTypeWind               PsType = 9
	PsTypeCommercialStorage  PsType = 12
)

var psTypeNames = map[PsType]string{
	PsTypeGround:             "地面电站",
	PsTypeDistributedPV:      "分布式光伏",
	PsTypeResidentialPV:      "户用光伏",
	PsTypeResidentialStorage: "户用储能",
	PsTypeVillage:            "村级电站",
	PsTypeDistributedStorage: "分布式储能",
	PsTypePovertyRelief:      "扶贫电站",
	PsTypeWind:               "风能电站",
	PsTypeCommercialStorage:  "工商业储能",
}

func (t PsType) String() string {
	if n, ok := psTypeNames[t]; ok {
		return n
	}
	return "未知类型电站"
}

type PsOnlineStatus uint8

const (
	PsOnlineStatusOffline PsOnlineStatus = 0
	PsOnlineStatusOnlie   PsOnlineStatus = 1
)

func (t PsOnlineStatus) String() string {
	switch t {
	case PsOnlineStatusOffline:
		return "离线"
	case PsOnlineStatusOnlie:
		return "在线"
	default:
		return "未知状态"
	}
}
