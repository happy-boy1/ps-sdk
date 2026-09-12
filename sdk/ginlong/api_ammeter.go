package ginlong

import (
	"fmt"
	"strings"
)

// AmmeterListRequest 账号下电表列表请求
type AmmeterListRequest struct {
	PageRequest

	StationID Int64 `json:"stationId,omitempty"` // 电站ID，空查全部
}

// AmmeterListItem 电表列表项
type AmmeterListItem struct {
	rawHolder

	ID            Int64       `json:"id"`            // 电表ID
	SN            string      `json:"sn"`            // 电表SN
	CollectorID   Int64       `json:"collectorId"`   // 采集器ID
	CollectorSN   string      `json:"collectorSn"`   // 采集器SN
	Name          string      `json:"name"`          // 电表名称
	StationID     Int64       `json:"stationId"`     // 电站ID
	StationName   string      `json:"stationName"`   // 电站名称
	State         DeviceState `json:"state"`         // 状态：1在线2离线
	DataTimestamp Int64       `json:"dataTimestamp"` // 更新时间戳(ms)
	Psum          Num         `json:"psum"`          // 实时功率
	PsumStr       string      `json:"psumStr"`       // 实时功率单位
	PsumPec       Num         `json:"psumPec"`       // 实时功率倍率

	ETodayPositiveActive    Num    `json:"eTodayPositiveActive"`    // 日正向电量
	ETodayPositiveActiveStr string `json:"eTodayPositiveActiveStr"` // 日正向电量单位
	ETodayReverseActive     Num    `json:"eTodayReverseActive"`     // 日反向电量
	ETodayReverseActiveStr  string `json:"eTodayReverseActiveStr"`  // 日反向电量单位
	ETotalPositiveActive    Num    `json:"eTotalPositiveActive"`    // 累计正向电量
	ETotalPositiveActiveStr string `json:"eTotalPositiveActiveStr"` // 累计正向电量单位
	ETotalReverseActive     Num    `json:"eTotalReverseActive"`     // 累计反向电量
	ETotalReverseActiveStr  string `json:"eTotalReverseActiveStr"`  // 累计反向电量单位

	// 以下字段未列入《返回参数》表，出现在文档响应示例中
	DataTimestampStr   string      `json:"dataTimestampStr"`   // 更新时间文本
	StateExceptionFlag int         `json:"stateExceptionFlag"` // 状态异常标志
	CollectorName      string      `json:"collectorName"`      // 采集器名称
	CollectorState     DeviceState `json:"collectorState"`     // 采集器状态
	AmmeterID          Int64       `json:"ammeterId"`          // 电表ID，同 id
	CurrentTransformer Num         `json:"currentTransformer"` // 电流互感器倍率
	TimeZone           int         `json:"timeZone"`           // 时区偏移(小时)
	TimeZoneStr        string      `json:"timeZoneStr"`        // 时区文本
	TimeZoneName       string      `json:"timeZoneName"`       // 时区名称
}

// AmmeterListResult 账号下电表列表结果
type AmmeterListResult struct {
	rawHolder

	Page Page[AmmeterListItem] `json:"page"` // 电表列表分页
}

// AmmeterList 获取账号下电表列表
func (c *SolisSDK) AmmeterList(req AmmeterListRequest) (*AmmeterListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()

	var out AmmeterListResult
	if err := c.do(PathAmmeterList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AmmeterDetailRequest 单台电表详情请求
type AmmeterDetailRequest struct {
	SN string `json:"sn"` // 电表SN，必填
}

// AmmeterDetailResult 单台电表详情结果
type AmmeterDetailResult struct {
	rawHolder

	ID                 Int64       `json:"id"`                 // 电表ID
	SN                 string      `json:"sn"`                 // 电表SN
	CollectorID        Int64       `json:"collectorId"`        // 采集器ID
	CollectorSN        string      `json:"collectorSn"`        // 采集器SN
	Name               string      `json:"name"`               // 电表名称
	StationID          Int64       `json:"stationId"`          // 电站ID
	StationName        string      `json:"stationName"`        // 电站名称
	State              DeviceState `json:"state"`              // 状态：1在线2离线
	DataTimestamp      Int64       `json:"dataTimestamp"`      // 更新时间戳(ms)
	CurrentTransformer Num         `json:"currentTransformer"` // 电流互感器倍率
	ElectricMeter      int         `json:"electricMeter"`      // 电表相数，0单相
	PA                 Num         `json:"pA"`                 // 有功功率A(W)
	PB                 Num         `json:"pB"`                 // 有功功率B(W)
	PC                 Num         `json:"pC"`                 // 有功功率C(W)
	AReactivePower     Num         `json:"aReactivePower"`     // 无功功率A(Var)
	BReactivePower     Num         `json:"bReactivePower"`     // 无功功率B(Var)
	CReactivePower     Num         `json:"cReactivePower"`     // 无功功率C(Var)
	ALookedPower       Num         `json:"aLookedPower"`       // 视在功率A(VA)
	BLookedPower       Num         `json:"bLookedPower"`       // 视在功率B(VA)
	CLookedPower       Num         `json:"cLookedPower"`       // 视在功率C(VA)
	UA                 Num         `json:"uA"`                 // A相电压(V)
	UB                 Num         `json:"uB"`                 // B相电压(V)
	UC                 Num         `json:"uC"`                 // C相电压(V)
	IA                 Num         `json:"iA"`                 // A相电流(A)
	IB                 Num         `json:"iB"`                 // B相电流(A)
	IC                 Num         `json:"iC"`                 // C相电流(A)
	Psum               Num         `json:"psum"`               // 实时功率
	PsumStr            string      `json:"psumStr"`            // 实时功率单位
	PsumPec            Num         `json:"psumPec"`            // 实时功率倍率
	TotalReactivePower Num         `json:"totalReactivePower"` // 无功总功率(Var)
	TotalViewPower     Num         `json:"totalViewPower"`     // 视在总功率(VA)
	FAc                Num         `json:"fAc"`                // 频率(Hz)
	AveragePowerFactor Num         `json:"averagePowerFactor"` // 功率因数

	ETodayPositiveActive    Num    `json:"eTodayPositiveActive"`    // 日正向电量
	ETodayPositiveActiveStr string `json:"eTodayPositiveActiveStr"` // 日正向电量单位
	ETodayReverseActive     Num    `json:"eTodayReverseActive"`     // 日反向电量
	ETodayReverseActiveStr  string `json:"eTodayReverseActiveStr"`  // 日反向电量单位
	EMonthPositiveActive    Num    `json:"eMonthPositiveActive"`    // 月正向电量
	EMonthPositiveActiveStr string `json:"eMonthPositiveActiveStr"` // 月正向电量单位
	EMonthReverseActive     Num    `json:"eMonthReverseActive"`     // 月反向电量
	EMonthReverseActiveStr  string `json:"eMonthReverseActiveStr"`  // 月反向电量单位
	EYearPositiveActive     Num    `json:"eYearPositiveActive"`     // 年正向电量
	EYearPositiveActiveStr  string `json:"eYearPositiveActiveStr"`  // 年正向电量单位
	EYearReverseActive      Num    `json:"eYearReverseActive"`      // 年反向电量
	EYearReverseActiveStr   string `json:"eYearReverseActiveStr"`   // 年反向电量单位
	ETotalPositiveActive    Num    `json:"eTotalPositiveActive"`    // 累计正向电量
	ETotalPositiveActiveStr string `json:"eTotalPositiveActiveStr"` // 累计正向电量单位
	ETotalReverseActive     Num    `json:"eTotalReverseActive"`     // 累计反向电量
	ETotalReverseActiveStr  string `json:"eTotalReverseActiveStr"`  // 累计反向电量单位

	// 以下字段未列入《返回参数》表，出现在文档响应示例中
	StationType           StationType `json:"stationType"`           // 电站类型，详见附录2
	TimeZone              int         `json:"timeZone"`              // 时区偏移(小时)
	TimeZoneStr           string      `json:"timeZoneStr"`           // 时区文本
	Daylight              int         `json:"daylight"`              // 夏令时
	DaylightSwitch        int         `json:"daylightSwitch"`        // 夏令时开关，0关1开
	Sno                   string      `json:"sno"`                   // 电站短ID
	CollectorState        DeviceState `json:"collectorState"`        // 采集器状态
	SimFlowState          int         `json:"simFlowState"`          // 流量卡状态
	StateExceptionFlag    int         `json:"stateExceptionFlag"`    // 状态异常标志
	ElectricMeterProtocol int         `json:"electricMeterProtocol"` // 电表协议类型
	IAStr                 string      `json:"iAStr"`                 // A相电流单位
	UAStr                 string      `json:"uAStr"`                 // A相电压单位
	IBStr                 string      `json:"iBStr"`                 // B相电流单位
	UBStr                 string      `json:"uBStr"`                 // B相电压单位
	ICStr                 string      `json:"iCStr"`                 // C相电流单位
	UCStr                 string      `json:"uCStr"`                 // C相电压单位
	TotalReactivePowerStr string      `json:"totalReactivePowerStr"` // 无功总功率单位
	TotalViewPowerStr     string      `json:"totalViewPowerStr"`     // 视在总功率单位
	PAStr                 string      `json:"pAStr"`                 // 有功功率A单位
	AReactivePowerStr     string      `json:"aReactivePowerStr"`     // 无功功率A单位
	ALookedPowerStr       string      `json:"aLookedPowerStr"`       // 视在功率A单位
	PBStr                 string      `json:"pBStr"`                 // 有功功率B单位
	BReactivePowerStr     string      `json:"bReactivePowerStr"`     // 无功功率B单位
	BLookedPowerStr       string      `json:"bLookedPowerStr"`       // 视在功率B单位
	PCStr                 string      `json:"pCStr"`                 // 有功功率C单位
	CReactivePowerStr     string      `json:"cReactivePowerStr"`     // 无功功率C单位
	CLookedPowerStr       string      `json:"cLookedPowerStr"`       // 视在功率C单位
	AphasePowerFactor     Num         `json:"aphasePowerFactor"`     // A相功率因数
	BphasePowerFactor     Num         `json:"bphasePowerFactor"`     // B相功率因数
	CphasePowerFactor     Num         `json:"cphasePowerFactor"`     // C相功率因数
}

// AmmeterDetail 获取单台电表详情
func (c *SolisSDK) AmmeterDetail(req AmmeterDetailRequest) (*AmmeterDetailResult, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("sn 不能为空")
	}

	var out AmmeterDetailResult
	if err := c.do(PathAmmeterDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
