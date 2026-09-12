package huawei

import (
	"fmt"
	"strings"
)

// StationKpiRequest 电站报表请求（小时/日/月/年）
type StationKpiRequest struct {
	StationCodes string `json:"stationCodes"` // 电站编号，英文逗号分隔，最多 100 个
	CollectTime  int64  `json:"collectTime"`  // 采集时间戳(ms)，决定统计的时间范围
}

// StationKpiItem 电站报表数据
type StationKpiItem struct {
	StationCode string  `json:"stationCode"` // 电站编号
	CollectTime int64   `json:"collectTime"` // 采集时间戳(ms)
	DataItemMap ItemMap `json:"dataItemMap"` // 数据项
}

// DeviceKpiRequest 设备报表请求（日/月/年）
type DeviceKpiRequest struct {
	DevIDs      string    `json:"devIds,omitempty"` // 设备编号，逗号分隔，与 Sns 二选一
	Sns         string    `json:"sns,omitempty"`    // 设备 SN，逗号分隔，与 DevIDs 二选一
	DevTypeID   DevTypeID `json:"devTypeId"`        // 设备类型，一次仅支持一种
	CollectTime int64     `json:"collectTime"`      // 采集时间戳(ms)，决定统计的时间范围
}

// DeviceKpiItem 设备报表数据
type DeviceKpiItem struct {
	DevID       int64   `json:"devId"`       // 设备编号
	Sn          string  `json:"sn"`          // 设备 SN
	CollectTime int64   `json:"collectTime"` // 采集时间戳(ms)
	DataItemMap ItemMap `json:"dataItemMap"` // 数据项
}

// GetStationKpiHour 查询电站小时级数据（collectTime 所在自然日的逐小时数据）
func (sdk *FusionSolarSDK) GetStationKpiHour(req StationKpiRequest) ([]StationKpiItem, error) {
	return sdk.stationKpi(PathGetKpiStationHour, req)
}

// GetStationKpiDay 查询电站日数据（collectTime 所在自然月的逐日数据）
func (sdk *FusionSolarSDK) GetStationKpiDay(req StationKpiRequest) ([]StationKpiItem, error) {
	return sdk.stationKpi(PathGetKpiStationDay, req)
}

// GetStationKpiMonth 查询电站月数据（collectTime 所在自然年的逐月数据）
func (sdk *FusionSolarSDK) GetStationKpiMonth(req StationKpiRequest) ([]StationKpiItem, error) {
	return sdk.stationKpi(PathGetKpiStationMonth, req)
}

// GetStationKpiYear 查询电站年数据
func (sdk *FusionSolarSDK) GetStationKpiYear(req StationKpiRequest) ([]StationKpiItem, error) {
	return sdk.stationKpi(PathGetKpiStationYear, req)
}

func (sdk *FusionSolarSDK) stationKpi(path string, req StationKpiRequest) ([]StationKpiItem, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	var out []StationKpiItem
	if err := sdk.do(path, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r StationKpiRequest) validate() error {
	if strings.TrimSpace(r.StationCodes) == "" {
		return fmt.Errorf("stationCodes 不能为空")
	}
	if r.CollectTime <= 0 {
		return fmt.Errorf("collectTime 不能为空")
	}
	if n := len(SplitCodes(r.StationCodes)); n > MaxBatch {
		return fmt.Errorf("stationCodes 最多 %d 个电站，当前 %d 个", MaxBatch, n)
	}
	return nil
}

// GetDevKpiDay 查询设备日数据（collectTime 所在自然月的逐日数据）
func (sdk *FusionSolarSDK) GetDevKpiDay(req DeviceKpiRequest) ([]DeviceKpiItem, error) {
	return sdk.deviceKpi(PathGetDevKpiDay, req)
}

// GetDevKpiMonth 查询设备月数据（collectTime 所在自然年的逐月数据）
func (sdk *FusionSolarSDK) GetDevKpiMonth(req DeviceKpiRequest) ([]DeviceKpiItem, error) {
	return sdk.deviceKpi(PathGetDevKpiMonth, req)
}

// GetDevKpiYear 查询设备年数据
func (sdk *FusionSolarSDK) GetDevKpiYear(req DeviceKpiRequest) ([]DeviceKpiItem, error) {
	return sdk.deviceKpi(PathGetDevKpiYear, req)
}

func (sdk *FusionSolarSDK) deviceKpi(path string, req DeviceKpiRequest) ([]DeviceKpiItem, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	var out []DeviceKpiItem
	if err := sdk.do(path, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r DeviceKpiRequest) validate() error {
	if r.DevTypeID == 0 {
		return fmt.Errorf("devTypeId 不能为空")
	}
	if strings.TrimSpace(r.DevIDs) == "" && strings.TrimSpace(r.Sns) == "" {
		return fmt.Errorf("devIds 与 sns 必须二选一")
	}
	if r.CollectTime <= 0 {
		return fmt.Errorf("collectTime 不能为空")
	}
	return nil
}
