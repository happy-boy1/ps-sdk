package huawei

import (
	"fmt"
	"strings"
)

// maxBatteryTaskPlants 单次任务最多下发的电站数
const maxBatteryTaskPlants = 10

// ==================== 5.2.3 储能工作模式设置任务下发 ====================

// BatteryModeTaskRequest 储能工作模式设置任务下发请求体
type BatteryModeTaskRequest struct {
	Tasks []BatteryModeTaskItem `json:"tasks"` // 任务列表，最多10个电站
}

// BatteryModeTaskItem 单个电站的储能工作模式设置任务
type BatteryModeTaskItem struct {
	PlantCode                        string                      `json:"plantCode"`                                  // 电站DN
	OperationMode                    BatteryOperationMode        `json:"operationMode"`                              // 工作模式
	RedundantPVEnergyPriority        PVEnergyPriority            `json:"redundantPVEnergyPriority,omitempty"`        // 多余PV能量优先级
	AllowedAcChargePower             *Float64                    `json:"allowedAcChargePower,omitempty"`             // 电网充电最大功率kW
	ChargingAndDischargingTimeWindow []ChargeDischargeTimeWindow `json:"chargingAndDischargingTimeWindow,omitempty"` // 充放电时间窗口
}

// ChargeDischargeTimeWindow 充放电时间窗口，最多14个
type ChargeDischargeTimeWindow struct {
	StartTime         string            `json:"startTime"`         // 开始时间，格式HH:MM
	EndTime           string            `json:"endTime"`           // 结束时间，格式HH:MM
	ChargeOrDischarge ChargeOrDischarge `json:"chargeOrDischarge"` // 充电/放电
	Repeat            []int             `json:"repeat"`            // 重复日期，1周一~7周日
}

// BatteryModeTaskResult 储能工作模式设置任务下发返回数据
type BatteryModeTaskResult struct {
	TaskID string                      `json:"taskId"` // 任务唯一ID
	Result []BatteryModeTaskResultItem `json:"result"` // 各电站下发结果
}

// BatteryModeTaskResultItem 单个电站的任务下发结果
type BatteryModeTaskResultItem struct {
	PlantCode string     `json:"plantCode"`         // 电站DN
	Status    TaskStatus `json:"status"`            // 下发任务当前状态
	Message   string     `json:"message,omitempty"` // 下发结果描述
}

// BatteryModeTask 下发储能工作模式设置任务
func (sdk *FusionSolarSDK) BatteryModeTask(req BatteryModeTaskRequest) (*BatteryModeTaskResult, error) {
	if len(req.Tasks) == 0 {
		return nil, fmt.Errorf("tasks 不能为空")
	}
	if len(req.Tasks) > maxBatteryTaskPlants {
		return nil, fmt.Errorf("tasks 数量不能超过 %d 个", maxBatteryTaskPlants)
	}
	for i, task := range req.Tasks {
		if strings.TrimSpace(task.PlantCode) == "" {
			return nil, fmt.Errorf("tasks[%d].plantCode 不能为空", i)
		}
		if strings.TrimSpace(string(task.OperationMode)) == "" {
			return nil, fmt.Errorf("tasks[%d].operationMode 不能为空", i)
		}
		if task.OperationMode == ModeTOU && len(task.ChargingAndDischargingTimeWindow) == 0 {
			return nil, fmt.Errorf("tasks[%d] 为TOU模式时 chargingAndDischargingTimeWindow 不能为空", i)
		}
		for j, w := range task.ChargingAndDischargingTimeWindow {
			if strings.TrimSpace(w.StartTime) == "" || strings.TrimSpace(w.EndTime) == "" {
				return nil, fmt.Errorf("tasks[%d].chargingAndDischargingTimeWindow[%d] 的 startTime/endTime 不能为空", i, j)
			}
			if strings.TrimSpace(string(w.ChargeOrDischarge)) == "" {
				return nil, fmt.Errorf("tasks[%d].chargingAndDischargingTimeWindow[%d].chargeOrDischarge 不能为空", i, j)
			}
			if len(w.Repeat) == 0 {
				return nil, fmt.Errorf("tasks[%d].chargingAndDischargingTimeWindow[%d].repeat 不能为空", i, j)
			}
		}
	}

	var result BatteryModeTaskResult
	if err := sdk.doTask(PathBatteryModeTask, req, &result); err != nil {
		return &result, err
	}
	return &result, nil
}

// ==================== 5.2.4 储能工作模式设置任务查询 ====================

// BatteryModeTaskInfoRequest 储能工作模式设置任务查询请求体
type BatteryModeTaskInfoRequest struct {
	TaskID string `json:"taskId"` // 任务唯一ID，16位
}

// BatteryModeTaskInfoResult 储能工作模式设置任务查询返回数据
type BatteryModeTaskInfoResult struct {
	DispatchResult []BatteryModeDispatchResult `json:"dispatchResult"` // 各电站下发结果
	StartTime      Time                        `json:"startTime"`      // 收到任务的时间
	EndTime        Time                        `json:"endTime"`        // 任务完成时间，未完成为空
}

// BatteryModeDispatchResult 单个电站的工作模式任务执行结果
type BatteryModeDispatchResult struct {
	PlantCode                        string                      `json:"plantCode"`                                  // 电站DN
	Status                           TaskStatus                  `json:"status"`                                     // 下发任务当前状态
	Message                          string                      `json:"message,omitempty"`                          // 失败原因，其他为空
	OperationMode                    BatteryOperationMode        `json:"operationMode"`                              // 工作模式
	RedundantPVEnergyPriority        PVEnergyPriority            `json:"redundantPVEnergyPriority,omitempty"`        // 多余PV能量优先级
	AllowedAcChargePower             *Float64                    `json:"allowedAcChargePower,omitempty"`             // 电网充电最大功率kW
	ChargingAndDischargingTimeWindow []ChargeDischargeTimeWindow `json:"chargingAndDischargingTimeWindow,omitempty"` // 充放电时间窗口
}

// BatteryModeTaskInfo 查询储能工作模式设置任务执行情况
func (sdk *FusionSolarSDK) BatteryModeTaskInfo(req BatteryModeTaskInfoRequest) (*BatteryModeTaskInfoResult, error) {
	if strings.TrimSpace(req.TaskID) == "" {
		return nil, fmt.Errorf("taskId 不能为空")
	}

	var result BatteryModeTaskInfoResult
	if err := sdk.do(PathBatteryModeTaskInfo, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ==================== 5.2.5 储能参数设置任务下发 ====================

// BatteryConfigTaskRequest 储能参数设置任务下发请求体
type BatteryConfigTaskRequest struct {
	Tasks []BatteryConfigTaskItem `json:"tasks"` // 任务列表，最多10个电站
}

// BatteryConfigTaskItem 单个电站的储能参数设置任务
type BatteryConfigTaskItem struct {
	PlantCode                string                   `json:"plantCode"`                // 电站DN
	BatteryConfigurationInfo BatteryConfigurationInfo `json:"batteryConfigurationInfo"` // 储能参数设置有效参数
}

// BatteryConfigurationInfo 储能参数设置有效参数，四项至少填一项
type BatteryConfigurationInfo struct {
	EndOfChargeSoc        *Float64 `json:"endOfChargeSoc,omitempty"`        // 充电截止SOC，单位%
	EndOfDischargeSoc     *Float64 `json:"endOfDischargeSoc,omitempty"`     // 放电截止SOC，单位%
	MaximumChargePower    *Float64 `json:"maximumChargePower,omitempty"`    // 最大充电功率，单位W
	MaximumDischargePower *Float64 `json:"maximumDischargePower,omitempty"` // 最大放电功率，单位W
}

// BatteryConfigTaskResult 储能参数设置任务下发返回数据
type BatteryConfigTaskResult struct {
	TaskID string                        `json:"taskId"` // 任务唯一ID
	Result []BatteryConfigTaskResultItem `json:"result"` // 各电站下发结果
}

// BatteryConfigTaskResultItem 单个电站的任务下发结果
type BatteryConfigTaskResultItem struct {
	PlantCode string     `json:"plantCode"`         // 电站DN
	Status    TaskStatus `json:"status"`            // 下发任务当前状态
	Message   string     `json:"message,omitempty"` // 下发结果描述
}

// BatteryConfigTask 下发储能参数设置任务
func (sdk *FusionSolarSDK) BatteryConfigTask(req BatteryConfigTaskRequest) (*BatteryConfigTaskResult, error) {
	if len(req.Tasks) == 0 {
		return nil, fmt.Errorf("tasks 不能为空")
	}
	if len(req.Tasks) > maxBatteryTaskPlants {
		return nil, fmt.Errorf("tasks 数量不能超过 %d 个", maxBatteryTaskPlants)
	}
	for i, task := range req.Tasks {
		if strings.TrimSpace(task.PlantCode) == "" {
			return nil, fmt.Errorf("tasks[%d].plantCode 不能为空", i)
		}
		info := task.BatteryConfigurationInfo
		if info.EndOfChargeSoc == nil && info.EndOfDischargeSoc == nil &&
			info.MaximumChargePower == nil && info.MaximumDischargePower == nil {
			return nil, fmt.Errorf("tasks[%d].batteryConfigurationInfo 至少需要填一个参数", i)
		}
	}

	var result BatteryConfigTaskResult
	if err := sdk.doTask(PathBatteryConfigTask, req, &result); err != nil {
		return &result, err
	}
	return &result, nil
}

// ==================== 5.2.6 储能参数设置任务查询 ====================

// BatteryConfigTaskInfoRequest 储能参数设置任务查询请求体
type BatteryConfigTaskInfoRequest struct {
	TaskID string `json:"taskId"` // 任务唯一ID，16位
}

// BatteryConfigTaskInfoResult 储能参数设置任务查询返回数据
type BatteryConfigTaskInfoResult struct {
	DispatchResult []BatteryConfigDispatchResult `json:"dispatchResult"` // 各电站下发结果
	StartTime      Time                          `json:"startTime"`      // 收到任务的时间
	EndTime        Time                          `json:"endTime"`        // 任务完成时间，未完成为空
}

// BatteryConfigDispatchResult 单个电站的参数设置任务执行结果
type BatteryConfigDispatchResult struct {
	PlantCode                string                   `json:"plantCode"`                // 电站DN
	Status                   TaskStatus               `json:"status"`                   // 下发任务当前状态
	Message                  string                   `json:"message,omitempty"`        // 失败原因，其他为空
	BatteryConfigurationInfo BatteryConfigurationInfo `json:"batteryConfigurationInfo"` // 储能参数设置有效参数
}

// BatteryConfigTaskInfo 查询储能参数设置任务执行情况
func (sdk *FusionSolarSDK) BatteryConfigTaskInfo(req BatteryConfigTaskInfoRequest) (*BatteryConfigTaskInfoResult, error) {
	if strings.TrimSpace(req.TaskID) == "" {
		return nil, fmt.Errorf("taskId 不能为空")
	}

	var result BatteryConfigTaskInfoResult
	if err := sdk.do(PathBatteryConfigTaskInfo, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
