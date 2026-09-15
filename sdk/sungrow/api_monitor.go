package sungrow

import (
	"encoding/json"
	"errors"
	"regexp"
)

type PageRequest struct {
	CurPage int `json:"curPage"`
	Size    int `json:"size,omitempty"`
}

type PageResult[T any] struct {
	PageList []T `json:"pageList"`
	RowCount int `json:"rowCount"`
}

// 电站列表信息查询
type PowerStationListRequest struct {
	Request

	PageRequest
	PsName    string `json:"ps_name,omitempty"`
	PsType    string `json:"ps_type,omitempty"`
	ShareType string `json:"share_type,omitempty"`
	ValidFlag string `json:"valid_flag,omitempty"`
	OrgID     string `json:"org_id,omitempty"`
}

type PowerStation struct {
	AlarmCount               int64          `json:"alarm_count"`
	BuildStatus              int64          `json:"build_status"`
	CityName                 string         `json:"city_name"`
	Co2Reduce                Entity         `json:"co2_reduce"`
	Co2ReduceTotal           Entity         `json:"co2_reduce_total"`
	Co2ReduceTotalUpdateTime string         `json:"co2_reduce_total_update_time"`
	Co2ReduceUpdateTime      string         `json:"co2_reduce_update_time"`
	ConnectType              int64          `json:"connect_type"`
	CurrPower                Entity         `json:"curr_power"`
	CurrPowerUpdateTime      string         `json:"curr_power_update_time"`
	Description              string         `json:"description"`
	DistrictName             string         `json:"district_name"`
	EquivalentHour           Entity         `json:"equivalent_hour"`
	EquivalentHourUpdateTime string         `json:"equivalent_hour_update_time"`
	FaultCount               int64          `json:"fault_count"`
	GridConnectionStatus     int64          `json:"grid_connection_status"`
	GridConnectionTime       int64          `json:"grid_connection_time"`
	InstallDate              string         `json:"install_date"`
	Latitude                 Num            `json:"latitude"`
	Longitude                Num            `json:"longitude"`
	MonthIncome              Entity         `json:"month_income"`
	MonthIncomeUpdateTime    string         `json:"month_income_update_time"`
	ProvinceName             string         `json:"province_name"`
	PsCurrentTimeZone        string         `json:"ps_current_time_zone"`
	PsFaultStatus            int64          `json:"ps_fault_status"`
	PsId                     Str            `json:"ps_id"`
	PsLocation               string         `json:"ps_location"`
	PsName                   string         `json:"ps_name"`
	PsStatus                 PsOnlineStatus `json:"ps_status"`
	PsType                   PsType         `json:"ps_type"`
	ShareType                string         `json:"share_type"`
	TodayEnergy              Entity         `json:"today_energy"`
	TodayEnergyUpdateTime    string         `json:"today_energy_update_time"`
	TodayIncome              Entity         `json:"today_income"`
	TodayIncomeUpdateTime    string         `json:"today_income_update_time"`
	TotalCapcity             Entity         `json:"total_capcity"`
	TotalCapcityUpdateTime   string         `json:"total_capcity_update_time"`
	TotalEnergy              Entity         `json:"total_energy"`
	TotalEnergyUpdateTime    string         `json:"total_energy_update_time"`
	TotalIncome              Entity         `json:"total_income"`
	TotalIncomeUpdateTime    string         `json:"total_income_update_time"`
	ValidFlag                int64          `json:"valid_flag"`
	YearIncome               Entity         `json:"year_income"`
	YearIncomeUpdateTime     string         `json:"year_income_update_time"`
}

type PowerStationList PageResult[PowerStation]

func (sdk *SungrowSDK) GetPowerStationList(req PowerStationListRequest) (*PowerStationList, error) {
	if req.CurPage < 0 {
		return nil, errors.New("curPage 参数必须大于0")
	}

	var out PowerStationList
	if err := sdk.do(PathGetPsList, &req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// 查询电站下设备列表

type DeviceListByPsIDRequest struct {
	Request

	PageRequest
	PsID                 string      `json:"ps_id"`
	IsVirtualUnit        string      `json:"is_virtual_unit,omitempty"`
	DeviceTypeList       []DevTypeID `json:"device_type_list,omitempty"`
	ClaimState           string      `json:"claim_state,omitempty"`
	IsGetFirmwareVersion string      `json:"is_get_firmware_version,omitempty"`
}

type Device struct {
	ChnnlId            int64     `json:"chnnl_id"`
	CommunicationDevSn string    `json:"communication_dev_sn"`
	DevFaultStatus     int64     `json:"dev_fault_status"`
	DevStatus          string    `json:"dev_status"`
	DeviceCode         int64     `json:"device_code"`
	DeviceModelCode    string    `json:"device_model_code"`
	DeviceModelId      int64     `json:"device_model_id"`
	DeviceName         string    `json:"device_name"`
	DeviceSn           string    `json:"device_sn"`
	DeviceType         DevTypeID `json:"device_type"`
	FactoryName        string    `json:"factory_name"`
	GridConnectionDate string    `json:"grid_connection_date"`
	PsId               int64     `json:"ps_id"`
	PsKey              string    `json:"ps_key"`
	RelState           int64     `json:"rel_state"`
	RelTime            string    `json:"rel_time"`
	TypeName           string    `json:"type_name"`
	Uuid               int64     `json:"uuid"`
}

type DeviceList PageResult[Device]

func (sdk *SungrowSDK) GetDeviceListByPsID(req DeviceListByPsIDRequest) (*DeviceList, error) {
	if req.CurPage < 0 {
		return nil, errors.New("curPage 参数必须大于0")
	}

	var out DeviceList
	if err := sdk.do(PathGetDevList, &req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// 测点类数据返回结果
type PointDataResult[T any] struct {
	FailPsKeyList   []string    `json:"fail_ps_key_list"`
	DevicePointList []T         `json:"device_point_list"`
	PointDict       []PointDict `json:"point_dict,omitempty"`
}

type PointDict struct {
	PointId   int64  `json:"point_id"`
	PointUnit string `json:"point_unit"`
	PointName string `json:"point_name"`
}

type DevicePoint struct {
	DevicePoint DevicePointInner `json:"device_point"`
}

type DevicePointInner struct {
	PsKey              string         `json:"ps_key"`
	DeviceSn           string         `json:"device_sn"`
	DevStatus          int64          `json:"dev_status"`
	Uuid               int64          `json:"uuid"`
	DeviceName         string         `json:"device_name"`
	DevFaultStatus     int64          `json:"dev_fault_status"`
	PsId               int64          `json:"ps_id"`
	CommunicationDevSn string         `json:"communication_dev_sn"`
	DeviceTime         string         `json:"device_time"`
	Points             map[string]any `json:"-"`
}

var pointFieldRegex = regexp.MustCompile(`^p\d+$`)

func (d *DevicePointInner) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.PsKey, _ = raw["ps_key"].(string)
	d.DeviceSn, _ = raw["device_sn"].(string)
	if v, ok := raw["dev_status"].(float64); ok {
		d.DevStatus = int64(v)
	}
	if v, ok := raw["uuid"].(float64); ok {
		d.Uuid = int64(v)
	}
	d.DeviceName, _ = raw["device_name"].(string)
	if v, ok := raw["dev_fault_status"].(float64); ok {
		d.DevFaultStatus = int64(v)
	}
	if v, ok := raw["ps_id"].(float64); ok {
		d.PsId = int64(v)
	}
	d.CommunicationDevSn, _ = raw["communication_dev_sn"].(string)
	d.DeviceTime, _ = raw["device_time"].(string)

	d.Points = make(map[string]any)
	for k, v := range raw {
		if pointFieldRegex.MatchString(k) {
			d.Points[k] = v
		}
	}
	return nil
}

func (d *DevicePointInner) MarshalJSON() ([]byte, error) {
	out := map[string]any{
		"ps_key":               d.PsKey,
		"device_sn":            d.DeviceSn,
		"dev_status":           d.DevStatus,
		"uuid":                 d.Uuid,
		"device_name":          d.DeviceName,
		"dev_fault_status":     d.DevFaultStatus,
		"ps_id":                d.PsId,
		"communication_dev_sn": d.CommunicationDevSn,
		"device_time":          d.DeviceTime,
	}

	for k, v := range d.Points {
		out[k] = v
	}

	return json.Marshal(out)
}

// 查询设备实时测点数据
type DeviceRtdRequest struct {
	Request

	PsKeyList      []string  `json:"ps_key_list,omitempty"`
	SnList         []string  `json:"sn_list,omitempty"`
	PointIDList    []string  `json:"point_id_list"`
	DeviceType     DevTypeID `json:"device_type"`
	IsGetPointDict string    `json:"is_get_point_dict,omitempty"`
}

type DeviceRealTimeData PointDataResult[DevicePoint]

func (sdk *SungrowSDK) GetDeviceRealTimeData(req DeviceRtdRequest) (*DeviceRealTimeData, error) {
	var out DeviceRealTimeData
	if err := sdk.do(PathGetDevRtd, &req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
