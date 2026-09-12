package huawei

import (
	"fmt"
	"strings"
)

// BatteryDispatchTaskRequest 储能调度充放电任务下发请求
type BatteryDispatchTaskRequest struct {
	PlantCode     string `json:"plantCode"`     // 电站DN
	BatteryDn     string `json:"batteryDn"`     // 储能或光储一体机DN
	DispatchPower int    `json:"dispatchPower"` // 充放电功率(W)，充电>0 放电<0 不充不放=0
}

// BatteryDispatchTaskResult 储能调度充放电任务下发结果
type BatteryDispatchTaskResult struct {
	TaskID    string     `json:"taskId"`    // 任务唯一ID，用于查询结果
	PlantCode string     `json:"plantCode"` // 电站DN
	BatteryDn string     `json:"batteryDn"` // 储能DN
	Status    TaskStatus `json:"status"`    // 任务状态：RUNNING/FAIL
}

// BatteryDispatchTask 下发储能调度充放电任务，须先将储能工作模式设为三方调度
func (sdk *FusionSolarSDK) BatteryDispatchTask(req BatteryDispatchTaskRequest) (*BatteryDispatchTaskResult, error) {
	if strings.TrimSpace(req.PlantCode) == "" {
		return nil, fmt.Errorf("plantCode 不能为空")
	}
	if strings.TrimSpace(req.BatteryDn) == "" {
		return nil, fmt.Errorf("batteryDn 不能为空")
	}
	var out BatteryDispatchTaskResult
	if err := sdk.doTask(PathBatteryDispatchTask, req, &out); err != nil {
		return &out, err
	}
	return &out, nil
}

// BatteryDispatchTaskInfoRequest 储能调度充放电任务查询请求
type BatteryDispatchTaskInfoRequest struct {
	TaskID string `json:"taskId"` // 任务唯一ID，固定 16 位
}

// BatteryDispatchTaskInfoResult 储能调度充放电任务查询结果
type BatteryDispatchTaskInfoResult struct {
	PlantCode     string         `json:"plantCode"`     // 电站DN
	BatteryDn     string         `json:"batteryDn"`     // 储能或光储一体机DN
	Status        TaskStatus     `json:"status"`        // 任务状态：RUNNING/SUCCESS/FAIL
	Message       TaskFailReason `json:"message"`       // 失败原因，成功时不返回
	DispatchPower int            `json:"dispatchPower"` // 实际下发的充放电功率(W)
}

// BatteryDispatchTaskInfo 查询储能调度充放电任务执行情况
func (sdk *FusionSolarSDK) BatteryDispatchTaskInfo(req BatteryDispatchTaskInfoRequest) (*BatteryDispatchTaskInfoResult, error) {
	if strings.TrimSpace(req.TaskID) == "" {
		return nil, fmt.Errorf("taskId 不能为空")
	}
	var out BatteryDispatchTaskInfoResult
	if err := sdk.do(PathBatteryDispatchTaskInfo, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
