package ginlong

import (
	"fmt"
	"strings"
)

// CollectorListRequest 账号下采集器列表请求
type CollectorListRequest struct {
	PageRequest        // 分页参数，pageNo/pageSize 必填
	StationID   Int64  `json:"stationId,omitempty"` // 电站id，空为全部
	NMICode     string `json:"nmiCode,omitempty"`   // nmi码，空为全部
}

// CollectionStatusVo 采集器状态统计
type CollectionStatusVo struct {
	All     int `json:"all"`     // 采集器总数
	Normal  int `json:"normal"`  // 采集器正常数
	Offline int `json:"offline"` // 采集器离线数
	Fault   int `json:"fault"`   // 采集器故障数
}

// CollectorListItem 采集器列表项
type CollectorListItem struct {
	rawHolder
	ID            Int64       `json:"id"`            // 采集器id
	StationName   string      `json:"stationName"`   // 电站名称
	StationID     Int64       `json:"stationId"`     // 电站Id
	UserID        Int64       `json:"userId"`        // 业主Id
	SN            string      `json:"sn"`            // 采集器SN
	Model         string      `json:"model"`         // 采集器型号
	Name          string      `json:"name"`          // 采集器名称
	RssiLevel     int         `json:"rssiLevel"`     // 采集器信号强度
	State         DeviceState `json:"state"`         // 状态：1在线/2离线/3报警
	DataTimestamp Int64       `json:"dataTimestamp"` // 更新时间
	ContractTime  Int64       `json:"contractTime"`  // 流量到期时间
}

// CollectorListResult 账号下采集器列表结果
type CollectorListResult struct {
	rawHolder
	Page               Page[CollectorListItem] `json:"page"`               // 结果列表
	CollectionStatusVo CollectionStatusVo      `json:"collectionStatusVo"` // 采集器状态统计
}

// CollectorList 获取账号下采集器列表
func (c *SolisSDK) CollectorList(req CollectorListRequest) (*CollectorListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out CollectorListResult
	if err := c.do(PathCollectorList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CollectorDetailRequest 单台采集器详情请求
type CollectorDetailRequest struct {
	ID Int64  `json:"id,omitempty"` // 采集器id，与 SN 至少填一个
	SN string `json:"sn,omitempty"` // 采集器SN，与 id 至少填一个
}

// CollectorDetailResult 单台采集器详情结果
type CollectorDetailResult struct {
	rawHolder
	ID                 Int64       `json:"id"`                 // 采集器id
	StationID          Int64       `json:"stationId"`          // 电站Id
	StationName        string      `json:"stationName"`        // 电站名称
	Addr               string      `json:"addr"`               // 电站地址
	UserID             Int64       `json:"userId"`             // 业主Id
	State              DeviceState `json:"state"`              // 状态：1在线/2离线/3报警
	DataTimestamp      Int64       `json:"dataTimestamp"`      // 电站更新时间
	TotalWorkingTime   Int64       `json:"totalWorkingTime"`   // 累计工作时间
	SN                 string      `json:"sn"`                 // 采集器SN
	Model              string      `json:"model"`              // 采集器型号
	Name               string      `json:"name"`               // 采集器名称
	RssiLevel          int         `json:"rssiLevel"`          // 采集器信号强度
	Lac                string      `json:"lac"`                // 定位lac
	LanIP              string      `json:"lanIp"`              // 局域网ip
	Mac                string      `json:"mac"`                // Mac地址
	MaximumNumber      int         `json:"maximumNumber"`      // 最大连接台数
	ActualNumber       int         `json:"actualNumber"`       // 实际连接台数
	ConnectedSSID      string      `json:"connectedSsid"`      // 连接的ssid
	ConnectionOperator string      `json:"connectionOperator"` // 运营商
	CurrentWorkingTime Int64       `json:"currentWorkingTime"` // 本次上电工作时间
	DataUploadCycle    int         `json:"dataUploadCycle"`    // 数据上传间隔
	FactoryTime        Int64       `json:"factoryTime"`        // 出厂时间
	ContractTime       Int64       `json:"contractTime"`       // 流量到期时间
}

// CollectorDetail 获取单台采集器详情
func (c *SolisSDK) CollectorDetail(req CollectorDetailRequest) (*CollectorDetailResult, error) {
	if req.ID == 0 && strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("id 和 sn 不能同时为空")
	}
	var out CollectorDetailResult
	if err := c.do(PathCollectorDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CollectorDayRequest 单台采集器信号值请求
type CollectorDayRequest struct {
	SN       string `json:"sn"`       // 采集器SN
	Time     string `json:"time"`     // 查询日期，例：2022-02-02
	TimeZone *Num   `json:"timeZone"` // 采集器所在时区，例：8
}

// CollectorDayItem 采集器当日信号值
type CollectorDayItem struct {
	rawHolder
	CollectorID   Int64  `json:"collectorId"`   // 采集器ID
	CollectorSN   string `json:"collectorSn"`   // 采集器SN
	DataTimestamp Int64  `json:"dataTimestamp"` // 更新时间(8时区)
	TimeStr       string `json:"timeStr"`       // 采集器时区的时间
	Pec           int    `json:"pec"`           // 信号强度百分比，单位%
	Rssi          int    `json:"rssi"`          // 信号强度值
	RssiLevel     int    `json:"rssiLevel"`     // 信号强度等级
}

// CollectorDay 获取单台采集器信号值
func (c *SolisSDK) CollectorDay(req CollectorDayRequest) ([]CollectorDayItem, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("sn 不能为空")
	}
	if strings.TrimSpace(req.Time) == "" {
		return nil, fmt.Errorf("time 不能为空")
	}
	if req.TimeZone == nil {
		return nil, fmt.Errorf("timeZone 不能为空")
	}
	var out []CollectorDayItem
	if err := c.do(PathCollectorDay, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeatherListRequest 账号下气象仪列表请求
type WeatherListRequest struct {
	PageRequest        // 分页参数，pageNo/pageSize 必填
	StationID   Int64  `json:"stationId,omitempty"` // 电站id，空为全部
	NMICode     string `json:"nmiCode,omitempty"`   // nmi码，空为全部
}

// WeatherListItem 气象仪列表项
type WeatherListItem struct {
	rawHolder
	ID              Int64       `json:"id"`              // 气象仪ID
	CollectorSN     string      `json:"collectorSn"`     // 气象仪/采集器SN
	CollectorID     Int64       `json:"collectorId"`     // 采集器ID
	Name            string      `json:"name"`            // 气象仪名称
	WeatherModel    string      `json:"weatherModel"`    // 型号：1阳光/2利诚
	UserID          Int64       `json:"userId"`          // 业主id
	StationID       Int64       `json:"stationId"`       // 电站id
	StationName     string      `json:"stationName"`     // 电站名称
	State           DeviceState `json:"state"`           // 状态：1在线/2离线
	DataTimestamp   Int64       `json:"dataTimestamp"`   // 更新时间
	TotalR          Num         `json:"totalR"`          // 总辐射，单位W/㎡
	DirectR         Num         `json:"directR"`         // 直接辐射，单位W/㎡
	ScatteredR      Num         `json:"scatteredR"`      // 散射辐射，单位W/㎡
	SunshineTim     Num         `json:"sunshineTim"`     // 日照时长，单位m
	TotalRDay       Num         `json:"totalRday"`       // 总辐射日累计，MJ/㎡
	DirectRDay      Num         `json:"directRday"`      // 直接辐射日累计，MJ/㎡
	ScatteredRDay   Num         `json:"scatteredRday"`   // 散射辐射日累计，MJ/㎡
	Temp            Num         `json:"temp"`            // 温度，单位取温度单位
	TemperatureUnit string      `json:"temperatureUnit"` // 温度单位
	Humidity        Num         `json:"humidity"`        // 湿度，单位%RH
	WindDirection   Num         `json:"windDirection"`   // 风向
	WindSpeed       Num         `json:"windSpeed"`       // 风速，单位m/s
	AirPressure     Num         `json:"airPressure"`     // 气压，单位Pa
	Rainfall        Num         `json:"rainfall"`        // 雨量，单位mm
	PvTemp          Num         `json:"pvTemp"`          // 组件温度，单位同上
}

// WeatherListResult 账号下气象仪列表结果
type WeatherListResult struct {
	rawHolder
	Page Page[WeatherListItem] `json:"page"` // 结果列表
}

// WeatherList 获取账号下气象仪列表
func (c *SolisSDK) WeatherList(req WeatherListRequest) (*WeatherListResult, error) {
	req.PageRequest = req.PageRequest.Normalize()
	var out WeatherListResult
	if err := c.do(PathWeatherList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WeatherDetailRequest 单台气象仪详情请求
type WeatherDetailRequest struct {
	SN string `json:"sn"` // 气象仪SN
}

// WeatherDetailResult 单台气象仪详情结果
type WeatherDetailResult struct {
	rawHolder
	ID              Int64       `json:"id"`              // 气象仪ID
	CollectorSN     string      `json:"collectorSn"`     // 气象仪/采集器SN
	CollectorID     Int64       `json:"collectorId"`     // 采集器ID
	Name            string      `json:"name"`            // 气象仪名称
	WeatherModel    string      `json:"weatherModel"`    // 型号：1阳光/2利诚
	UserID          Int64       `json:"userId"`          // 业主id
	StationID       Int64       `json:"stationId"`       // 电站id
	StationName     string      `json:"stationName"`     // 电站名称
	State           DeviceState `json:"state"`           // 状态：1在线/2离线
	DataTimestamp   Int64       `json:"dataTimestamp"`   // 更新时间
	TotalR          Num         `json:"totalR"`          // 总辐射，单位W/㎡
	DirectR         Num         `json:"directR"`         // 直接辐射，单位W/㎡
	ScatteredR      Num         `json:"scatteredR"`      // 散射辐射，单位W/㎡
	SunshineTim     Num         `json:"sunshineTim"`     // 日照时长，单位m
	TotalRDay       Num         `json:"totalRday"`       // 总辐射日累计，MJ/㎡
	DirectRDay      Num         `json:"directRday"`      // 直接辐射日累计，MJ/㎡
	ScatteredRDay   Num         `json:"scatteredRday"`   // 散射辐射日累计，MJ/㎡
	Temp            Num         `json:"temp"`            // 温度，单位取温度单位
	TemperatureUnit string      `json:"temperatureUnit"` // 温度单位
	Humidity        Num         `json:"humidity"`        // 湿度，单位%RH
	WindDirection   Num         `json:"windDirection"`   // 风向
	WindSpeed       Num         `json:"windSpeed"`       // 风速，单位m/s
	AirPressure     Num         `json:"airPressure"`     // 气压，单位Pa
	Rainfall        Num         `json:"rainfall"`        // 雨量，单位mm
	PvTemp          Num         `json:"pvTemp"`          // 组件温度，单位同上
}

// WeatherDetail 获取单台气象仪详情
func (c *SolisSDK) WeatherDetail(req WeatherDetailRequest) (*WeatherDetailResult, error) {
	if strings.TrimSpace(req.SN) == "" {
		return nil, fmt.Errorf("sn 不能为空")
	}
	var out WeatherDetailResult
	if err := c.do(PathWeatherDetail, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
