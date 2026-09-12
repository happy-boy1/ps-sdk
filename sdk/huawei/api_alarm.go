package huawei

import (
	"fmt"
	"strings"
)

// AlarmListRequest 活动告警查询请求
type AlarmListRequest struct {
	StationCodes string `json:"stationCodes,omitempty"` // 电站编号，与 Sns 至少填一个
	Sns          string `json:"sns,omitempty"`          // 设备 SN，与 StationCodes 至少填一个
	BeginTime    int64  `json:"beginTime"`              // 开始时间戳(ms)
	EndTime      int64  `json:"endTime"`                // 结束时间戳(ms)
	Language     string `json:"language"`               // 语言，如 zh_CN、en_US
	Levels       string `json:"levels,omitempty"`       // 告警级别，如 "1,2"，空为全部
	DevTypes     string `json:"devTypes,omitempty"`     // 设备类型，如 "1,38"，空为全部
}

// Alarm 活动告警
type Alarm struct {
	StationCode      string    `json:"stationCode"`      // 电站编号
	StationName      string    `json:"stationName"`      // 电站名称
	AlarmID          int       `json:"alarmId"`          // 告警 ID
	AlarmName        string    `json:"alarmName"`        // 告警名称
	AlarmCause       string    `json:"alarmCause"`       // 告警原因
	AlarmType        int       `json:"alarmType"`        // 告警类型
	CauseID          int       `json:"causeId"`          // 原因 ID
	DevName          string    `json:"devName"`          // 设备名称
	DevTypeID        DevTypeID `json:"devTypeId"`        // 设备类型
	EsnCode          string    `json:"esnCode"`          // 设备 SN
	RepairSuggestion string    `json:"repairSuggestion"` // 修复建议
	RaiseTime        int64     `json:"raiseTime"`        // 告警产生时间戳(ms)
	Lev              int       `json:"lev"`              // 告警级别 1紧急 2重要 3次要 4提示
	Status           int       `json:"status"`           // 告警状态 1未处理
}

// GetAlarmList 查询当前活动告警（不含历史告警）
func (sdk *FusionSolarSDK) GetAlarmList(req AlarmListRequest) ([]Alarm, error) {
	if strings.TrimSpace(req.StationCodes) == "" && strings.TrimSpace(req.Sns) == "" {
		return nil, fmt.Errorf("stationCodes 与 sns 至少填一个")
	}
	if req.BeginTime <= 0 || req.EndTime <= 0 {
		return nil, fmt.Errorf("beginTime/endTime 不能为空")
	}
	if strings.TrimSpace(req.Language) == "" {
		return nil, fmt.Errorf("language 不能为空")
	}
	var out []Alarm
	if err := sdk.do(PathGetAlarmList, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
