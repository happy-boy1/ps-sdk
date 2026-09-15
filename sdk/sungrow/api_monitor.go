package sungrow

import (
	"errors"
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
	AlarmCount               int64   `json:"alarm_count"`
	BuildStatus              int64   `json:"build_status"`
	CityName                 string  `json:"city_name"`
	Co2Reduce                Entity  `json:"co2_reduce"`
	Co2ReduceTotal           Entity  `json:"co2_reduce_total"`
	Co2ReduceTotalUpdateTime string  `json:"co2_reduce_total_update_time"`
	Co2ReduceUpdateTime      string  `json:"co2_reduce_update_time"`
	ConnectType              int64   `json:"connect_type"`
	CurrPower                Entity  `json:"curr_power"`
	CurrPowerUpdateTime      string  `json:"curr_power_update_time"`
	Description              string  `json:"description"`
	DistrictName             string  `json:"district_name"`
	EquivalentHour           Entity  `json:"equivalent_hour"`
	EquivalentHourUpdateTime string  `json:"equivalent_hour_update_time"`
	FaultCount               int64   `json:"fault_count"`
	GridConnectionStatus     int64   `json:"grid_connection_status"`
	GridConnectionTime       int64   `json:"grid_connection_time"`
	InstallDate              string  `json:"install_date"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
	MonthIncome              Entity  `json:"month_income"`
	MonthIncomeUpdateTime    string  `json:"month_income_update_time"`
	ProvinceName             string  `json:"province_name"`
	PsCurrentTimeZone        string  `json:"ps_current_time_zone"`
	PsFaultStatus            int64   `json:"ps_fault_status"`
	PsId                     int64   `json:"ps_id"`
	PsLocation               string  `json:"ps_location"`
	PsName                   string  `json:"ps_name"`
	PsStatus                 int64   `json:"ps_status"`
	PsType                   int64   `json:"ps_type"`
	ShareType                string  `json:"share_type"`
	TodayEnergy              Entity  `json:"today_energy"`
	TodayEnergyUpdateTime    string  `json:"today_energy_update_time"`
	TodayIncome              Entity  `json:"today_income"`
	TodayIncomeUpdateTime    string  `json:"today_income_update_time"`
	TotalCapcity             Entity  `json:"total_capcity"`
	TotalCapcityUpdateTime   string  `json:"total_capcity_update_time"`
	TotalEnergy              Entity  `json:"total_energy"`
	TotalEnergyUpdateTime    string  `json:"total_energy_update_time"`
	TotalIncome              Entity  `json:"total_income"`
	TotalIncomeUpdateTime    string  `json:"total_income_update_time"`
	ValidFlag                int64   `json:"valid_flag"`
	YearIncome               Entity  `json:"year_income"`
	YearIncomeUpdateTime     string  `json:"year_income_update_time"`
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
