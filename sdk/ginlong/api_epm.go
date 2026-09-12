package ginlong

import (
	"fmt"
	"strings"
)

// EpmListRequest 获取账号下 EPM 列表请求
type EpmListRequest struct {
	PageRequest
	StationID string `json:"stationId"` // 电站ID，空为账号下全部
	NmiCode   string `json:"nmiCode"`   // nmi编码，空为账号下全部
}

// EpmListItem EPM 列表条目
type EpmListItem struct {
	rawHolder
	ID            Int64       `json:"id"`            // EPM_ID
	SN            string      `json:"sn"`            // EPM_SN
	CollectorID   Int64       `json:"collectorId"`   // 采集器ID
	CollectorSN   string      `json:"collectorSn"`   // 采集器SN
	UserID        Int64       `json:"userId"`        // 业主ID
	StationID     Int64       `json:"stationId"`     // 电站ID
	StationName   string      `json:"stationName"`   // 电站名称
	State         DeviceState `json:"state"`         // 状态：1在线2离线
	DataTimestamp Int64       `json:"dataTimestamp"` // 更新时间
	FailSafe      int         `json:"failSafe"`      // failSafe开关
	PEpmTotal     Num         `json:"pEpmTotal"`     // EPM总功率
	PEpmTotalStr  string      `json:"pEpmTotalStr"`  // EPM总功率单位
	ETotalBuy     Num         `json:"eTotalBuy"`     // 总买电
	ETotalBuyStr  string      `json:"eTotalBuyStr"`  // 总买电单位
	ETotalSell    Num         `json:"eTotalSell"`    // 总卖电
	ETotalSellStr string      `json:"eTotalSellStr"` // 总卖电单位
}

// EpmListResult 获取账号下 EPM 列表结果
type EpmListResult struct {
	rawHolder
	Page Page[EpmListItem] `json:"page"` // 结果列表
}

// EpmList 获取账号下 EPM 列表
func (c *SolisSDK) EpmList(req EpmListRequest) (*EpmListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out EpmListResult
	if err := c.do(PathEpmList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EpmDetailRequest 获取单台 EPM 详情请求
type EpmDetailRequest struct {
	SN string `json:"sn"` // EPM_SN，指定查询某台
}

// EpmDetailResult 获取单台 EPM 详情结果
type EpmDetailResult struct {
	rawHolder
	ID                 Int64       `json:"id"`                 // EPM_ID
	SN                 string      `json:"sn"`                 // EPM_SN
	CollectorID        Int64       `json:"collectorId"`        // 采集器ID
	CollectorSN        string      `json:"collectorSn"`        // 采集器SN
	UserID             Int64       `json:"userId"`             // 业主ID
	StationID          Int64       `json:"stationId"`          // 电站ID
	StationName        string      `json:"stationName"`        // 电站名称
	State              DeviceState `json:"state"`              // 状态：1在线2离线
	DataTimestamp      Int64       `json:"dataTimestamp"`      // 更新时间
	FailSafe           int         `json:"failSafe"`           // failSafe开关
	EmpSoftwareVersion Str         `json:"empSoftwareVersion"` // EPM软件版本，示例为数字
	PLimit             Num         `json:"pLimit"`             // 功率限制百分比
	CtRatio            Num         `json:"ctRatio"`            // 电流传感器变比
	PSet               Num         `json:"pSet"`               // 回流功率设置值
	PSetStr            string      `json:"pSetStr"`            // 回流功率设置值单位
	PInverterTotal     Num         `json:"pInverterTotal"`     // 逆变器总功率
	PInverterTotalStr  string      `json:"pInverterTotalStr"`  // 逆变器总功率单位
	EToaalInverter     Num         `json:"eToaalInverter"`     // 逆变器总发电量
	EToaalInverterStr  string      `json:"eToaalInverterStr"`  // 逆变器总发电量单位
	PLoad              Num         `json:"pLoad"`              // 用电总功率
	PLoadStr           string      `json:"pLoadStr"`           // 用电总功率单位
	ETotalLoad         Num         `json:"eTotalLoad"`         // 总用电量
	ETotalLoadStr      string      `json:"eTotalLoadStr"`      // 总用电量单位
	PEpmTotal          Num         `json:"pEpmTotal"`          // EPM总功率
	PEpmTotalStr       string      `json:"pEpmTotalStr"`       // EPM总功率单位
	ETotalBuy          Num         `json:"eTotalBuy"`          // 总买电
	ETotalBuyStr       string      `json:"eTotalBuyStr"`       // 总买电单位
	ETotalSell         Num         `json:"eTotalSell"`         // 总卖电
	ETotalSellStr      string      `json:"eTotalSellStr"`      // 总卖电单位
	IAc1               Num         `json:"iAc1"`               // 电流U
	IAc2               Num         `json:"iAc2"`               // 电流V
	IAc3               Num         `json:"iAc3"`               // 电流W
	UAc1               Num         `json:"uAc1"`               // 电压U
	UAc2               Num         `json:"uAc2"`               // 电压V
	UAc3               Num         `json:"uAc3"`               // 电压W
	PAc1               Num         `json:"pAc1"`               // 功率U
	PAc2               Num         `json:"pAc2"`               // 功率V
	PAc3               Num         `json:"pAc3"`               // 功率W
	PowerFactor        Num         `json:"powerFactor"`        // 功率因数
	FacMeter           Num         `json:"facMeter"`           // 电网频率
}

// EpmDetail 获取单台 EPM 详情
func (c *SolisSDK) EpmDetail(req EpmDetailRequest) (*EpmDetailResult, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathEpmDetail)
	}
	var out EpmDetailResult
	if err := c.do(PathEpmDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EpmDayRequest 获取单台 EPM 某日实时数据请求
type EpmDayRequest struct {
	SN         string `json:"sn"`         // EPM_SN
	SearchInfo string `json:"searchinfo"` // 查询字段，逗号分隔
	Time       string `json:"time"`       // 查询某日，如2022-09-24
	TimeZone   Num    `json:"timeZone"`   // 电站时区
}

// EpmDayItem 单台 EPM 某日实时数据条目
type EpmDayItem struct {
	rawHolder
	DataTimestamp  Int64  `json:"dataTimestamp"`  // 更新时间(8时区)
	TimeStr        string `json:"timeStr"`        // 按电站时区转换的时间
	PEpmTotal      Num    `json:"pEpmTotal"`      // EPM总功率(电网总功率)
	PEpmTotalStr   string `json:"pEpmTotalStr"`   // EPM总功率单位
	PEpmTotalPec   Num    `json:"pEpmTotalPec"`   // EPM总功率百分比
	ETotalBuy      Num    `json:"eTotalBuy"`      // 电网取电总有功电能
	ETotalSell     Num    `json:"eTotalSell"`     // 电网送电总有功电能
	UAc1           Num    `json:"uAc1"`           // EPM交流电压U
	IAc1           Num    `json:"iAc1"`           // EPM交流电流U
	PAc1           Num    `json:"pAc1"`           // EPM有功功率U
	UAc2           Num    `json:"uAc2"`           // EPM交流电压V
	IAc2           Num    `json:"iAc2"`           // EPM交流电流V
	PAc2           Num    `json:"pAc2"`           // EPM有功功率V
	UAc3           Num    `json:"uAc3"`           // EPM交流电压W
	IAc3           Num    `json:"iAc3"`           // EPM交流电流W
	PAc3           Num    `json:"pAc3"`           // EPM有功功率W
	PInverterTotal Num    `json:"pInverterTotal"` // 逆变器总功率
	PLimit         Num    `json:"pLimit"`         // 功率限制百分比
	CtRatio        Num    `json:"ctRatio"`        // CT电流传感器变比
	PowerFactor    Num    `json:"powerFactor"`    // 电网功率因数
	FacMeter       Num    `json:"facMeter"`       // 电网频率
	PLoad          Num    `json:"pLoad"`          // 负载总功率
	EToaalInverter Num    `json:"eToaalInverter"` // 逆变器总发电量
	ETotalLoad     Num    `json:"eTotalLoad"`     // 负载总用电量
}

// EpmDay 获取单台 EPM 某日的实时数据（图表接口）。
// 该接口文档表格写 data 为对象数组、官方示例却是「指标名→并列数组」的对象，故用 Chart 兼容两种形态：
// 数组形态读 Items，对象形态读 Raw（键名与 searchinfo 一致，如 uAc1、eTotalBuy、dataTimestamp）
func (c *SolisSDK) EpmDay(req EpmDayRequest) (*Chart[EpmDayItem], error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathEpmDay)
	}
	if strings.TrimSpace(req.SearchInfo) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 searchinfo 不能为空", PathEpmDay)
	}
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 time 不能为空", PathEpmDay)
	}
	out := new(Chart[EpmDayItem])
	if err := c.do(PathEpmDay, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// EpmMonthRequest 获取单台 EPM 某月日数据请求
type EpmMonthRequest struct {
	SN    string `json:"sn"`    // EPM_SN
	Month string `json:"month"` // 查询某月，如2019-07
}

// EpmMonthItem 单台 EPM 某月日数据条目
type EpmMonthItem struct {
	rawHolder
	Date          Int64  `json:"date"`          // 更新时间戳
	DateStr       string `json:"dateStr"`       // 更新时间字符串
	Energy        Num    `json:"energy"`        // 发电量
	EpmSellEnergy Num    `json:"epmSellEnergy"` // 卖电量
	EpmBuyEnergy  Num    `json:"epmBuyEnergy"`  // 买电量
}

// EpmMonth 获取单台 EPM 某月的日数据
func (c *SolisSDK) EpmMonth(req EpmMonthRequest) ([]EpmMonthItem, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathEpmMonth)
	}
	if strings.TrimSpace(req.Month) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 month 不能为空", PathEpmMonth)
	}
	var out []EpmMonthItem
	if err := c.do(PathEpmMonth, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EpmYearRequest 获取单台 EPM 某年月数据请求
type EpmYearRequest struct {
	SN   string `json:"sn"`   // EPM_SN
	Year string `json:"year"` // 查询某年，如2019
}

// EpmYearItem 单台 EPM 某年月数据条目
type EpmYearItem struct {
	rawHolder
	Date          Int64  `json:"date"`          // 更新时间戳(+8时区)
	DateStr       string `json:"dateStr"`       // 更新时间字符串
	Energy        Num    `json:"energy"`        // 发电量
	EpmSellEnergy Num    `json:"epmSellEnergy"` // 卖电量
	EpmBuyEnergy  Num    `json:"epmBuyEnergy"`  // 买电量
}

// EpmYear 获取单台 EPM 某年的月数据
func (c *SolisSDK) EpmYear(req EpmYearRequest) ([]EpmYearItem, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathEpmYear)
	}
	if strings.TrimSpace(req.Year) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 year 不能为空", PathEpmYear)
	}
	var out []EpmYearItem
	if err := c.do(PathEpmYear, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EpmAllRequest 获取单台 EPM 年数据请求
type EpmAllRequest struct {
	SN string `json:"sn"` // EPM_SN
}

// EpmAllItem 单台 EPM 年数据条目
type EpmAllItem struct {
	rawHolder
	Year          int `json:"year"`          // 年
	Energy        Num `json:"energy"`        // 发电量
	EpmSellEnergy Num `json:"epmSellEnergy"` // 卖电量
	EpmBuyEnergy  Num `json:"epmBuyEnergy"`  // 买电量
}

// EpmAll 获取单台 EPM 的年数据
func (c *SolisSDK) EpmAll(req EpmAllRequest) ([]EpmAllItem, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("[Ginlong]%s 必填参数 sn 不能为空", PathEpmAll)
	}
	var out []EpmAllItem
	if err := c.do(PathEpmAll, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
