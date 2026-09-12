package sungrow

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
	ResultCodeSuccess                                                   = "1"
	ResultCodeError                                                     = "-1"
	ResultCodeErUnknownException                                        = "000"
	ResultCodeErMissingParameterAppkey                                  = "001"
	ResultCodeErMissingParameterToken                                   = "002"
	ResultCodeErMissingParameterSysCode                                 = "003"
	ResultCodeErInvalidAppkey                                           = "E00000"
	ResultCodeErApiServiceHasExpired                                    = "E00001"
	ResultCodeErParameterDecryptError                                   = "E00002"
	ResultCodeErTokenLoginInvalid                                       = "E00003"
	ResultCodeErMonthCallApiTimesUpperLimit                             = "E998"
	ResultCodeErHourCallApiTimesUpperLimit                              = "E999"
	ResultCodeErMissingParameter                                        = "009"
	ResultCodeErParameterValueInvalid                                   = "010"
	ResultCodeErSqlException                                            = "011"
	ResultCodeUnauthorizedAccess                                        = "E900"
	ResultCodeCallTooFrequently                                         = "E901"
	ResultCodeAbnormalNetworkEnvironment                                = "E903"
	ResultCodeRequestIsNotEncrypted                                     = "E902"
	ResultCodeMissingParameterInRequestHeaderXRandomSecretKey           = "E904"
	ResultCodeAesDecryptionException                                    = "E905"
	ResultCodeRsaDecryptionException                                    = "E906"
	ResultCodeAesRandomSecretKeyLengthMustBe16                          = "E907"
	ResultCodeMissingKeyParameterApiKeyParam                            = "E908"
	ResultCodeInvalidParameterFormatNonce32BitStringOfNumbersAndLetters = "E909"
	ResultCodeRepeatedRequest                                           = "E910"
	ResultCodeMissingParameterInRequestHeaderXAccessKey                 = "E911"
	ResultCodeIllegalXAccessKey                                         = "E912"
	ResultCodeExpiredRequest                                            = "E913"
	ResultCodeMismatchedAppkeyAndXAccessKey                             = "E914"
	ResultCodeLoginTooFrequently                                        = "E916"
	ResultCodePermitDeniedByWhiteIpAddress                              = "E918"
	ResultCodePermitDeniedByWhiteListUser                               = "E919"
	ResultCodeSystemNotFound                                            = "E994"
	ResultCodeRequestBodyTooLarge                                       = "E995"
	ResultCodeApiNotFound                                               = "E996"
	ResultCodeTransformBusinessResponseDataOccurError                   = "E997"
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
