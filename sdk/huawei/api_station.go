package huawei

import (
	"fmt"
	"strings"
)

// StationListRequest 电站列表请求
type StationListRequest struct {
	PageNo                 int   `json:"pageNo"`                           // 分页页码，从 1 开始
	GridConnectedStartTime int64 `json:"gridConnectedStartTime,omitempty"` // 并网开始时间戳(ms)
	GridConnectedEndTime   int64 `json:"gridConnectedEndTime,omitempty"`   // 并网结束时间戳(ms)
}

// Station 电站信息
type Station struct {
	PlantCode          string  `json:"plantCode"`          // 电站编号
	PlantName          string  `json:"plantName"`          // 电站名称
	PlantAddress       string  `json:"plantAddress"`       // 详细地址
	Longitude          Float64 `json:"longitude"`          // 经度
	Latitude           Float64 `json:"latitude"`           // 纬度
	Capacity           Float64 `json:"capacity"`           // 组串总容量(kWp)
	ContactPerson      string  `json:"contactPerson"`      // 联系人
	ContactMethod      string  `json:"contactMethod"`      // 联系方式
	GridConnectionDate string  `json:"gridConnectionDate"` // 并网时间(含时区)
}

// StationListResult 电站列表分页结果
type StationListResult = PageResult[Station]

// GetStationList 查询电站列表，每页最多 100 条
func (sdk *FusionSolarSDK) GetStationList(req StationListRequest) (*StationListResult, error) {
	if req.PageNo <= 0 {
		return nil, fmt.Errorf("pageNo 必须大于 0")
	}
	var out StationListResult
	if err := sdk.do(PathStations, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeviceListRequest 设备列表请求
type DeviceListRequest struct {
	StationCodes string `json:"stationCodes"` // 电站编号，英文逗号分隔，最多 100 个
}

// Device 设备信息
type Device struct {
	ID              int64     `json:"id"`              // 设备编号
	DevDn           string    `json:"devDn"`           // 设备唯一编号
	DevName         string    `json:"devName"`         // 设备名称
	StationCode     string    `json:"stationCode"`     // 所属电站编号
	EsnCode         string    `json:"esnCode"`         // 设备 SN
	DevTypeID       DevTypeID `json:"devTypeId"`       // 设备类型
	Model           string    `json:"model"`           // 设备型号
	SoftwareVersion string    `json:"softwareVersion"` // 软件版本号
	OptimizerNumber int       `json:"optimizerNumber"` // 优化器数量
	InvType         string    `json:"invType"`         // 机型(仅逆变器)
	Longitude       Float64   `json:"longitude"`       // 经度
	Latitude        Float64   `json:"latitude"`        // 纬度
}

// GetDevList 查询设备列表，一次最多 100 个电站
func (sdk *FusionSolarSDK) GetDevList(req DeviceListRequest) ([]Device, error) {
	if strings.TrimSpace(req.StationCodes) == "" {
		return nil, fmt.Errorf("stationCodes 不能为空")
	}
	if n := len(SplitCodes(req.StationCodes)); n > MaxBatch {
		return nil, fmt.Errorf("stationCodes 最多 %d 个电站，当前 %d 个", MaxBatch, n)
	}
	var out []Device
	if err := sdk.do(PathGetDevList, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
