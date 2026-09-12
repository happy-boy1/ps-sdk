package huawei

import (
	"fmt"
	"strings"
)

// StationRealKpiRequest 电站实时数据请求
type StationRealKpiRequest struct {
	StationCodes string `json:"stationCodes"` // 电站编号，英文逗号分隔，最多 100 个
}

// StationRealKpi 电站实时数据
type StationRealKpi struct {
	StationCode string  `json:"stationCode"` // 电站编号
	DataItemMap ItemMap `json:"dataItemMap"` // 实时数据项，key 见文档电站实时数据列表
}

// GetStationRealKpi 查询电站实时数据，一次最多 100 个电站
func (sdk *FusionSolarSDK) GetStationRealKpi(req StationRealKpiRequest) ([]StationRealKpi, error) {
	if strings.TrimSpace(req.StationCodes) == "" {
		return nil, fmt.Errorf("stationCodes 不能为空")
	}
	var out []StationRealKpi
	if err := sdk.do(PathGetStationRealKpi, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceRealKpiRequest 设备实时数据请求
type DeviceRealKpiRequest struct {
	DevIDs    string    `json:"devIds,omitempty"` // 设备编号，逗号分隔，与 Sns 二选一
	Sns       string    `json:"sns,omitempty"`    // 设备 SN，逗号分隔，与 DevIDs 二选一
	DevTypeID DevTypeID `json:"devTypeId"`        // 设备类型，一次仅支持一种
}

// DeviceRealKpi 设备实时数据
type DeviceRealKpi struct {
	DevID       int64   `json:"devId"`       // 设备编号
	Sn          string  `json:"sn"`          // 设备 SN(储能返回 null)
	DataItemMap ItemMap `json:"dataItemMap"` // 实时数据项，N/A 表示无效值
}

// GetDevRealKpi 查询设备实时数据，一次最多 1 种设备类型 100 个设备
func (sdk *FusionSolarSDK) GetDevRealKpi(req DeviceRealKpiRequest) ([]DeviceRealKpi, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	var out []DeviceRealKpi
	if err := sdk.do(PathGetDevRealKpi, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r DeviceRealKpiRequest) validate() error {
	if r.DevTypeID == 0 {
		return fmt.Errorf("devTypeId 不能为空")
	}
	if strings.TrimSpace(r.DevIDs) == "" && strings.TrimSpace(r.Sns) == "" {
		return fmt.Errorf("devIds 与 sns 必须二选一")
	}
	return nil
}

// DeviceHistoryRequest 设备历史数据请求
type DeviceHistoryRequest struct {
	DevDn     string    `json:"devDn"`     // 设备编号，由设备列表 devDn 获取
	DevTypeID DevTypeID `json:"devTypeId"` // 设备类型
	StartTime int64     `json:"startTime"` // 开始时间戳(ms)
	EndTime   int64     `json:"endTime"`   // 结束时间戳(ms)
}

// DeviceHistoryItem 设备 5 分钟粒度历史数据
type DeviceHistoryItem struct {
	DevDns      string  `json:"devDns"`      // 设备编号
	CollectTime int64   `json:"collectTime"` // 采集时间戳(ms)
	DataItems   ItemMap `json:"dataItems"`   // 数据项，key 见文档设备历史数据列表
}

// GetDevHistory 查询设备历史数据，一次最多 1 种设备类型 1 个设备 24 小时
func (sdk *FusionSolarSDK) GetDevHistory(req DeviceHistoryRequest) ([]DeviceHistoryItem, error) {
	if strings.TrimSpace(req.DevDn) == "" {
		return nil, fmt.Errorf("devDn 不能为空")
	}
	if req.DevTypeID == 0 {
		return nil, fmt.Errorf("devTypeId 不能为空")
	}
	if req.StartTime <= 0 || req.EndTime <= 0 {
		return nil, fmt.Errorf("startTime/endTime 不能为空")
	}
	if req.StartTime >= req.EndTime {
		return nil, fmt.Errorf("startTime 必须小于 endTime")
	}
	var out []DeviceHistoryItem
	if err := sdk.do(PathDeviceHistory, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
